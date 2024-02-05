package ordercommand

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/app/queries/orderquery"
	"latipe-order-service-v2/internal/common/errors"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msgDTO"
	voucherConst "latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
	deliverygrpc "latipe-order-service-v2/internal/infrastructure/grpc/deliveryServ"
	productgrpc "latipe-order-service-v2/internal/infrastructure/grpc/productServ"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
	usergrpc "latipe-order-service-v2/internal/infrastructure/grpc/userServ"
	publishMsg "latipe-order-service-v2/internal/msgqueue"
	"latipe-order-service-v2/pkg/cache/redis"
	"latipe-order-service-v2/pkg/util/mapper"
	"strings"
)

type orderCommandService struct {
	cfg         *config.Config
	orderRepo   order.Repository
	cacheEngine *redis.CacheEngine
	publisher   *publishMsg.PublisherTransactionMessage
	//grpc client
	voucherGrpc  vouchergrpc.VoucherServiceClient
	productGrpc  productgrpc.ProductServiceClient
	deliveryGrpc deliverygrpc.DeliveryServiceClient
	userGrpc     usergrpc.UserServiceClient
}

func NewOrderCommmandService(cfg *config.Config, orderRepo order.Repository,
	cacheEngine *redis.CacheEngine, publisher *publishMsg.PublisherTransactionMessage,
	voucherGrpc vouchergrpc.VoucherServiceClient, productGrpc productgrpc.ProductServiceClient,
	deliveryGrpc deliverygrpc.DeliveryServiceClient, userGrpc usergrpc.UserServiceClient) OrderCommandUsecase {
	return orderCommandService{
		orderRepo:    orderRepo,
		cacheEngine:  cacheEngine,
		publisher:    publisher,
		cfg:          cfg,
		deliveryGrpc: deliveryGrpc,
		voucherGrpc:  voucherGrpc,
		productGrpc:  productGrpc,
		userGrpc:     userGrpc,
	}
}

func (o orderCommandService) CreateOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error) {
	userReq := usergrpc.GetDetailAddressRequest{AddressId: dto.Address.AddressId, UserId: dto.UserRequest.UserId}
	address, err := o.userGrpc.GetAddressDetail(ctx, &userReq)
	if err != nil {
		return nil, err
	}

	shopVoucherMap := make(map[string]string)
	cartMap := make(map[string][]string)
	//get product data from product service
	var orders []*order.Order

	// Handle order group by store
	checkout := msgDTO.CheckoutMessage{
		CheckoutID:    uuid.NewString(),
		UserID:        dto.UserRequest.UserId,
		PaymentMethod: dto.PaymentMethod,
	}

	for _, i := range dto.PromotionData.ShopVoucherInfo {
		shopVoucherMap[i.StoreId] = i.VoucherCode
	}

	for _, i := range dto.StoreOrders {

		productReq, itemMap := MappingToProductRequest(i)

		products, err := o.productGrpc.CheckInStock(ctx, productReq)
		if err != nil {
			return nil, err
		}

		//handle order group by store_order and calculate shipping cost
		shippingReq := deliverygrpc.GetShippingCostRequest{
			SrcCode:    products.ProvinceCode,
			DestCode:   address.CityOrProvinceId,
			DeliveryId: i.Delivery.DeliveryId,
		}
		shippingDetail, err := o.deliveryGrpc.CalculateShippingCost(ctx, &shippingReq)
		if err != nil {
			return nil, err
		}

		//init order data
		orderData, err := o.saveOrderIntoDatabase(ctx, dto, address, shippingDetail, products, itemMap, shopVoucherMap)
		if err != nil {
			log.Error(err)
			continue
		}

		//add key into response and send order data to orchestration service

		checkout.OrderData = append(checkout.OrderData, msgDTO.OrderData{
			OrderID: orderData.OrderID,
			Amount:  uint(orderData.Amount),
		})

		checkout.TotalAmount += uint(orderData.Amount)

		if len(i.CartIds) < 0 {
			cartMap[i.StoreID] = i.CartIds
		}

		orders = append(orders, orderData)
	}

	if err := o.publisher.SendOrderCreatedMessage(MappingDataToMessage(orders, cartMap)); err != nil {
		log.Error(err)
	}

	data := orderDTO.CreateOrderResponse{
		CheckoutMessage: checkout,
	}

	return &data, nil
}

func (o orderCommandService) saveOrderIntoDatabase(ctx context.Context, dto *orderDTO.CreateOrderRequest,
	address *usergrpc.GetDetailAddressResponse, deli *deliverygrpc.GetShippingCostResponse,
	productItems *productgrpc.GetPurchaseItemResponse, itemMap map[string]int, shopVoucherMap map[string]string) (*order.Order, error) {

	orderDAO := order.Order{}
	orderDAO.OrderID = GenKeyOrder(dto.UserRequest.UserId)
	orderDAO.Username = dto.UserRequest.Username
	orderDAO.UserId = dto.UserRequest.UserId
	orderDAO.StoreId = productItems.StoreId

	if err := mapper.BindingStruct(dto, &orderDAO); err != nil {
		log.Errorf("Mapping value from dto to dao failed cause: %s", err)
		return nil, err
	}

	orderDAO.ShippingCost = int(deli.Cost)

	var orderItems []*order.OrderItem
	for _, item := range productItems.Items {
		i := order.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.Name,
			NameOption:  item.NameOption,
			OptionID:    item.OptionId,
			Quantity:    GetQuantityItems(item.ProductId, item.OptionId, itemMap),
			Price:       int(item.Price),
			NetPrice:    int(item.PromotionalPrice),
			ProdImg:     item.Image,
			Order:       &orderDAO,
		}
		//calculate subtotal of item
		if i.NetPrice != 0 {
			i.SubTotal = i.NetPrice * i.Quantity
		} else {
			i.SubTotal = i.Price * i.Quantity
		}
		orderDAO.SubTotal += i.SubTotal
		orderItems = append(orderItems, &i)
	}
	orderDAO.OrderItem = orderItems

	// Check if the order has a voucher to apply a discount to the order
	if dto.PromotionData != nil {
		var voucherCode []string
		//shipping and payment vouchers
		voucherDetail, err := o.voucherGrpc.CheckUsingVouchersForCheckout(ctx,
			orderquery.MappingPaymentAndShippingVoucherRequest(dto, &orderDAO))
		if err != nil {
			return nil, err
		}

		if voucherDetail != nil {
			if voucherDetail.ShippingVoucher != nil {
				if uint64(deli.Cost) < voucherDetail.ShippingVoucher.DiscountData.ShippingValue {
					orderDAO.ShippingDiscount = int(deli.Cost)
				} else {
					orderDAO.ShippingDiscount = int(voucherDetail.ShippingVoucher.DiscountData.ShippingValue)
				}
				voucherCode = append(voucherCode, voucherDetail.ShippingVoucher.VoucherCode)
			}

			if voucherDetail.PaymentVoucher != nil {
				switch voucherDetail.PaymentVoucher.DiscountData.DiscountType {
				case voucherConst.FIXED_DISCOUNT:
					orderDAO.PaymentDiscount = int(voucherDetail.PaymentVoucher.DiscountData.DiscountValue)
				case voucherConst.PERCENT_DISCOUNT:
					value := uint64(float32(orderDAO.SubTotal) * voucherDetail.PaymentVoucher.DiscountData.DiscountPercent)

					if value <= voucherDetail.PaymentVoucher.DiscountData.MaximumValue {
						orderDAO.PaymentDiscount = int(value)
					} else {
						orderDAO.PaymentDiscount = int(voucherDetail.PaymentVoucher.DiscountData.MaximumValue)
					}
				}
				voucherCode = append(voucherCode, voucherDetail.PaymentVoucher.VoucherCode)
			}
		}

		//store vouchers
		if shopVoucherMap[orderDAO.StoreId] != "" {
			voucherResp, err := o.voucherGrpc.CheckUsingVouchersForCheckout(ctx,
				orderquery.MappingShopVoucherRequest(&orderDAO, shopVoucherMap[orderDAO.StoreId]))
			if err != nil {
				return nil, err
			}

			switch voucherResp.StoreVouchers[0].DiscountData.DiscountType {
			case voucherConst.FIXED_DISCOUNT:
				orderDAO.StoreDiscount += int(voucherResp.StoreVouchers[0].DiscountData.DiscountValue)
			case voucherConst.PERCENT_DISCOUNT:
				value := uint64(float32(orderDAO.SubTotal) * voucherResp.StoreVouchers[0].DiscountData.DiscountPercent)

				if value <= voucherResp.StoreVouchers[0].DiscountData.MaximumValue {
					orderDAO.StoreDiscount += int(value)
				} else {
					orderDAO.StoreDiscount += int(voucherResp.StoreVouchers[0].DiscountData.MaximumValue)
				}
			}
			voucherCode = append(voucherCode, shopVoucherMap[orderDAO.StoreId])

		}

		orderDAO.Vouchers = strings.Join(voucherCode, ";")
	}

	//calculate amount order
	orderDAO.Amount = (orderDAO.SubTotal + orderDAO.ShippingCost) -
		(orderDAO.PaymentDiscount + orderDAO.ShippingDiscount + orderDAO.StoreDiscount)
	if orderDAO.Amount < 0 {
		orderDAO.Amount = 0
	}
	orderDAO.Status = order.ORDER_SYSTEM_PROCESS

	//create delivery
	recvTime, err := order.ParseStringToDate(deli.ReceiveDate)
	if err != nil {
		return nil, err
	}

	shippingData := order.DeliveryOrder{
		DeliveryId:      deli.DeliveryId,
		DeliveryName:    deli.DeliveryName,
		Cost:            int(deli.Cost),
		AddressId:       address.Id,
		ShippingName:    address.ContactName,
		ShippingPhone:   address.Phone,
		ShippingAddress: address.DetailAddress,
		ReceivingDate:   *recvTime,
		Order:           &orderDAO,
	}
	orderDAO.Delivery = &shippingData

	orderDAO.PaymentMethod = dto.PaymentMethod
	orderDAO.Status = order.ORDER_SYSTEM_PROCESS
	//create log
	var logs []*order.OrderStatusLog
	orderLog := order.OrderStatusLog{
		Order:        &orderDAO,
		Message:      "Đơn hàng chờ hệ thống xử lý",
		StatusChange: orderDAO.Status,
	}
	orderDAO.OrderStatusLog = append(logs, &orderLog)

	//save order into db
	err = o.orderRepo.Save(ctx, &orderDAO)
	if err != nil {
		log.Errorf("the order created failed : %s", err)
		return nil, err
	}

	return &orderDAO, nil
}

func (o orderCommandService) StoreUpdateOrderStatus(ctx context.Context, dto *store.StoreUpdateOrderStatusRequest) error {
	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if orderDAO.StoreId != dto.StoreId {
		return errors.ErrNotFoundRecord
	}

	if orderDAO.Status != order.ORDER_CREATED {
		return errors.OrderStatusNotValid
	}

	switch dto.Status {
	case order.ORDER_PREPARED:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_PREPARED); err != nil {
			return err
		}
	case order.ORDER_CANCEL_BY_STORE:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_CANCEL_BY_STORE,
			fmt.Sprintf("Đơn hàng bị hủy do yêu cầu của cửa hàng:%v", dto.Message)); err != nil {
			return err
		}
	default:
		return errors.OrderCannotUpdate
	}

	return nil
}

func (o orderCommandService) DeliveryUpdateOrderStatus(ctx context.Context, dto delivery.UpdateOrderStatusRequest) error {
	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if orderDAO.Delivery.DeliveryId != dto.DeliveryID {
		return errors.ErrNotFoundRecord
	}

	if orderDAO.Status != order.ORDER_PREPARED {
		return errors.OrderStatusNotValid
	}

	switch dto.Status {
	case order.ORDER_DELIVERY:
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			"Đơn hàng được đơn vị vận chuyển tiếp nhận và giao hàng"); err != nil {
			return err
		}
	case order.ORDER_CANCEL_BY_USER:
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			fmt.Sprintf("Đơn hàng không thể giao thành công lý do: %v", dto.Message)); err != nil {
			return err
		}
	case order.ORDER_CANCEL_BY_DELI:
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			fmt.Sprintf("Đơn hàng giao thất bại lý do: %v", dto.Message)); err != nil {
			return err
		}
	case order.ORDER_SHIPPING_FINISH:
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			fmt.Sprintf("Đơn hàng giao thành công: %v", dto.Message)); err != nil {
			return err
		}
	default:
		return errors.OrderCannotUpdate
	}

	return nil
}

func (o orderCommandService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if orderDAO.Status < 0 {
		return errors.ErrBadRequest
	}

	orderDAO.Status = dto.Status

	if err := o.orderRepo.Update(ctx, *orderDAO); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) UpdateOrderStatusByReplyMessage(ctx context.Context, dto *msgDTO.OrderStatusMessage) error {
	switch dto.Status {
	case msgDTO.ORDER_EVENT_COMMIT_SUCCESS:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_CREATED,
			"Đơn hàng được tạo thành công"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PRODUCT:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi sản phẩm"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PROMOTION:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi khuyến mãi"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_DELIVERY:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi vận chuyển"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PAYMENT:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi thanh toán"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_CANCEL:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			"Đơn hàng bị hủy"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_REFUND:
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			"Đơn hàng hoàn trả"); err != nil {
			return err
		}
	}

	return nil
}

func (o orderCommandService) UserRefundOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if dao.UserId != dto.UserId {
		return errors.ErrNotFoundRecord
	}

	switch dao.Status {
	case order.ORDER_SHIPPING_FINISH:
		if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_REFUND); err != nil {
			return err
		}
	default:
		return errors.OrderCannotRefund
	}

	mess := msgDTO.OrderStatusMessage{
		OrderID: dao.OrderID,
		Status:  order.ORDER_REFUND,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) UserCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if dao.UserId != dto.UserId {
		return errors.ErrNotFoundRecord
	}

	switch dao.Status {
	case order.ORDER_CREATED:
		if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_CANCEL_BY_USER,
			fmt.Sprintf("Đơn hàng bị hủy do yêu cầu của người mua:%v", dto.Message)); err != nil {
			return err
		}
	default:
		return errors.OrderCannotCancel
	}

	mess := msgDTO.OrderStatusMessage{
		OrderID: dao.OrderID,
		Status:  order.ORDER_CANCEL_BY_USER,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_CANCEL_BY_ADMIN,
		fmt.Sprintf("Đơn hàng bị hủy do yêu cầu của quản trị viên:%v", dto.Message)); err != nil {
		return err
	}

	mess := msgDTO.OrderStatusMessage{
		OrderID: dao.OrderID,
		Status:  order.ORDER_CANCEL_BY_ADMIN,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}
