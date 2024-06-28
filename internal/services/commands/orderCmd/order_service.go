package orderCmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/common/errors"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msgDTO"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv"
	dto2 "latipe-order-service-v2/internal/infrastructure/adapter/storeserv/dto"
	voucherConst "latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
	deliverygrpc "latipe-order-service-v2/internal/infrastructure/grpc/deliveryServ"
	productgrpc "latipe-order-service-v2/internal/infrastructure/grpc/productServ"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
	usergrpc "latipe-order-service-v2/internal/infrastructure/grpc/userServ"
	publishMsg "latipe-order-service-v2/internal/publisher"
	"latipe-order-service-v2/internal/services/queries/orderQuery"
	cacheV9 "latipe-order-service-v2/pkg/cache/redisCacheV9"
	"latipe-order-service-v2/pkg/util/mapper"
	"strings"
	"time"
)

type orderCommandService struct {
	cfg            *config.Config
	orderRepo      order.OrderRepository
	commissionRepo order.CommissionRepository
	cacheEngine    *cacheV9.CacheEngine
	publisher      *publishMsg.PublisherTransactionMessage
	notifyPub      *publishMsg.NotificationMessagePublisher
	//grpc client
	voucherGrpc  vouchergrpc.VoucherServiceClient
	productGrpc  productgrpc.ProductServiceClient
	deliveryGrpc deliverygrpc.DeliveryServiceClient
	userGrpc     usergrpc.UserServiceClient
	//rest call
	storeServ storeserv.Service
}

func NewOrderCommandService(cfg *config.Config, orderRepo order.OrderRepository,
	commissionRepo order.CommissionRepository,
	cacheEngine *cacheV9.CacheEngine, publisher *publishMsg.PublisherTransactionMessage,
	voucherGrpc vouchergrpc.VoucherServiceClient, productGrpc productgrpc.ProductServiceClient,
	deliveryGrpc deliverygrpc.DeliveryServiceClient, userGrpc usergrpc.UserServiceClient, notifyPub *publishMsg.NotificationMessagePublisher,
	storeServ storeserv.Service) OrderCommandUsecase {
	return orderCommandService{
		orderRepo:      orderRepo,
		commissionRepo: commissionRepo,
		cacheEngine:    cacheEngine,
		publisher:      publisher,
		cfg:            cfg,
		deliveryGrpc:   deliveryGrpc,
		voucherGrpc:    voucherGrpc,
		productGrpc:    productGrpc,
		userGrpc:       userGrpc,
		storeServ:      storeServ,
		notifyPub:      notifyPub,
	}
}

func (o orderCommandService) CreateOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error) {
	userReq := usergrpc.GetDetailAddressRequest{AddressId: dto.Address.AddressId, UserId: dto.UserRequest.UserId}
	address, err := o.userGrpc.GetAddressDetail(ctx, &userReq)
	if err != nil {
		return nil, err
	}

	cartMap := make(map[string][]string)
	var orders []*order.Order
	var totalOrdersAmount int64
	// Handle order group by store
	checkout := msgDTO.CheckoutMessage{
		CheckoutID:    strings.ReplaceAll(uuid.NewString(), "-", ""),
		UserID:        dto.UserRequest.UserId,
		PaymentMethod: dto.PaymentMethod,
	}

	data := orderDTO.CreateOrderResponse{
		CheckoutMessage: checkout,
	}

	for _, i := range dto.StoreOrders {
		//mapping cart req data
		if len(i.CartIds) > 0 {
			cartMap[i.StoreID] = i.CartIds
		}

		//get product data from product service
		productReq, itemMap := MappingToProductRequest(i)
		products, err := o.productGrpc.CheckInStock(ctx, productReq)
		if err != nil {
			log.Error(err)
			data.FailedOrder.StoreID = i.StoreID
			data.FailedOrder.Message = "out of stock"
			continue
		}

		//handle order group by store_order and calculate shipping cost
		shippingReq := deliverygrpc.GetShippingCostRequest{
			SrcCode:    products.ProvinceCode,
			DestCode:   address.CityOrProvinceId,
			DeliveryId: i.Delivery.DeliveryId,
		}
		shippingDetail, err := o.deliveryGrpc.CalculateShippingCost(ctx, &shippingReq)
		if err != nil {
			data.FailedOrder.StoreID = i.StoreID
			data.FailedOrder.Message = "calculating shipping cost is failed"
			log.Error(err)
			continue
		}

		//init order data
		orderData, err := o.initOrderData(dto, address, shippingDetail, products, itemMap)
		if err != nil {
			data.FailedOrder.StoreID = i.StoreID
			data.FailedOrder.Message = err.Error()
			log.Error(err)
			continue
		}

		totalOrdersAmount += int64(orderData.SubTotal)
		orders = append(orders, orderData)
	}
	if len(orders) == 0 {
		return nil, errors.OrderCannotCreated
	}

	//handle promotion data
	if dto.PromotionData != nil {
		err = o.handlePromotionData(ctx, orders, dto.PromotionData, totalOrdersAmount)
		if err != nil {
			return nil, errors.OrderCannotCreated
		}
	}

	for _, entity := range orders {
		//add key into response and send order data to orchestration service
		checkout.OrderData = append(checkout.OrderData, msgDTO.OrderData{
			OrderID: entity.OrderID,
			Amount:  uint(entity.Amount),
		})
		checkout.TotalAmount += uint(entity.Amount)
	}

	if err := o.orderRepo.SaveMultiple(ctx, orders); err != nil {
		log.Error(err)
		return nil, errors.OrderCannotCreated
	}

	if err := o.publisher.SendOrderCreatedMessage(MappingDataToMessage(orders, cartMap, checkout)); err != nil {
		log.Error(err)
	}

	if err := o.cacheEngine.Set(checkout.CheckoutID, checkout, 24*time.Hour); err != nil {
		log.Error(err)
	}

	data.CheckoutMessage = checkout
	return &data, nil
}

func (o orderCommandService) initOrderData(dto *orderDTO.CreateOrderRequest,
	address *usergrpc.GetDetailAddressResponse, deli *deliverygrpc.GetShippingCostResponse,
	productItems *productgrpc.GetPurchaseItemResponse, itemMap map[string]int) (*order.Order, error) {

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

	//calculate amount order
	orderDAO.Amount = orderDAO.SubTotal + orderDAO.ShippingCost
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
	orderDAO.Thumbnail = orderItems[0].ProdImg
	//create log
	var logs []*order.OrderStatusLog
	orderLog := order.OrderStatusLog{
		Order:        &orderDAO,
		Message:      "Đơn hàng chờ hệ thống xử lý",
		StatusChange: orderDAO.Status,
	}
	orderDAO.OrderStatusLog = append(logs, &orderLog)

	return &orderDAO, nil
}

func (o orderCommandService) handlePromotionData(ctx context.Context, orders []*order.Order,
	promotionData *orderDTO.PromotionData, totalOrderAmount int64) error {

	var pVoucherResp *vouchergrpc.CheckoutVoucherResponse
	var err error

	if promotionData.PaymentVoucherInfo != nil {
		pVoucherReq := orderQuery.MappingPaymentVoucherRequest(promotionData, totalOrderAmount,
			orders[0].PaymentMethod, orders[0].UserId)
		pVoucherResp, err = o.voucherGrpc.CheckoutVoucherForPurchase(ctx, pVoucherReq)
		if err != nil {
			return err
		}
	}

	for _, entity := range orders {
		var voucherCode []string
		//shipping vouchers
		if promotionData.FreeShippingVoucherInfo != nil {
			sVoucherReq := orderQuery.MappingShippingVoucherRequest(promotionData, entity)
			if sVoucherReq == nil {
				continue
			}

			resp, err := o.voucherGrpc.CheckoutVoucherForPurchase(ctx, sVoucherReq)
			if err != nil {
				return err
			}

			if int64(entity.Amount) >= resp.VoucherDetail.VoucherRequire.MinRequire &&
				isMatchPaymentMethod(int(resp.VoucherDetail.VoucherRequire.PaymentMethod), entity.PaymentMethod) {

				if uint64(entity.ShippingCost) < resp.VoucherDetail.DiscountData.ShippingValue {
					entity.ShippingDiscount = entity.ShippingCost
				} else {
					entity.ShippingDiscount = int(resp.VoucherDetail.DiscountData.ShippingValue)
				}
				voucherCode = append(voucherCode, resp.VoucherDetail.VoucherCode)
			}
		}

		//store vouchers
		if len(promotionData.ShopVoucherInfo) > 0 {
			shopVoucherIndex := slices.IndexFunc(promotionData.ShopVoucherInfo, func(i orderDTO.ShopVoucherInfo) bool {
				return i.StoreId == entity.StoreId
			})
			if shopVoucherIndex != -1 {

				sVoucher := orderQuery.MappingShopVoucherRequest(entity,
					promotionData.ShopVoucherInfo[shopVoucherIndex].VoucherCode)

				voucherResp, err := o.voucherGrpc.CheckoutVoucherForPurchase(ctx, sVoucher)
				if err != nil {
					return err
				}

				switch voucherResp.VoucherDetail.DiscountData.DiscountType {

				case voucherConst.FIXED_DISCOUNT:
					entity.StoreDiscount += int(voucherResp.VoucherDetail.DiscountData.DiscountValue)
				case voucherConst.PERCENT_DISCOUNT:
					value := uint64(float32(entity.SubTotal) * voucherResp.VoucherDetail.DiscountData.DiscountPercent)

					if value <= voucherResp.VoucherDetail.DiscountData.MaximumValue {
						entity.StoreDiscount += int(value)
					} else {
						entity.StoreDiscount += int(voucherResp.VoucherDetail.DiscountData.MaximumValue)
					}
				}
				voucherCode = append(voucherCode, promotionData.ShopVoucherInfo[shopVoucherIndex].VoucherCode)

			}
		}

		if pVoucherResp != nil {
			switch pVoucherResp.VoucherDetail.DiscountData.DiscountType {
			case voucherConst.FIXED_DISCOUNT:
				entity.PaymentDiscount = int(pVoucherResp.VoucherDetail.DiscountData.DiscountValue) / len(orders)
			case voucherConst.PERCENT_DISCOUNT:
				value := uint64(float32(entity.SubTotal) * pVoucherResp.VoucherDetail.DiscountData.DiscountPercent)

				if value <= pVoucherResp.VoucherDetail.DiscountData.MaximumValue {
					entity.PaymentDiscount = int(value) / len(orders)
				} else {
					entity.PaymentDiscount = int(pVoucherResp.VoucherDetail.DiscountData.MaximumValue) / len(orders)
				}
			}
			voucherCode = append(voucherCode, pVoucherResp.VoucherDetail.VoucherCode)
		}

		entity.Amount = entity.SubTotal - (entity.PaymentDiscount + entity.StoreDiscount + entity.ShippingDiscount)
		entity.Vouchers = strings.Join(voucherCode, "-")
	}

	return nil
}

func (o orderCommandService) StoreUpdateOrderStatus(ctx context.Context, dto *store.StoreUpdateOrderStatusRequest) error {
	orderDAO, err := o.orderRepo.FindByIdForUpdate(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if orderDAO.StoreId != dto.StoreId {
		return errors.ErrNotFoundRecord
	}

	if orderDAO.Status != order.ORDER_CREATED {
		return errors.OrderStatusNotValid
	}

	var bodyContent string

	switch dto.Status {
	case order.ORDER_PREPARED:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_PREPARED); err != nil {
			return err
		}

		bodyContent = fmt.Sprintf("Xin chào, đơn hàng của bạn đã được cửa hàng xác nhận #[%v]", orderDAO.OrderID)
	case order.ORDER_CANCEL_BY_STORE:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_CANCEL_BY_STORE,
			fmt.Sprintf("Đơn hàng bị hủy do yêu cầu của cửa hàng:%v", dto.Message)); err != nil {
			return err
		}
		bodyContent = fmt.Sprintf("Xin lỗi đơn hàng #[%v] của bạn vừa bị hủy lý do là: %v ", orderDAO.OrderID,
			fmt.Sprintf("Đơn hàng bị hủy do yêu cầu của cửa hàng:%v", dto.Message))
	default:
		return errors.OrderCannotUpdate
	}

	userNoti := msgDTO.NewNotificationMessage(msgDTO.NOTIFY_USER, orderDAO.UserId, "[Latipe] Thông báo tình trạng đơn hàng", bodyContent, orderDAO.Thumbnail)

	//send message to user
	if err := o.notifyPub.NotifyToUser(userNoti); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) DeliveryUpdateOrderStatus(ctx context.Context, dto delivery.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindByIdForUpdate(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	var bodyContent string

	if orderDAO.Delivery.DeliveryId != dto.DeliveryID {
		return errors.ErrNotFoundRecord
	}

	if orderDAO.Status < order.ORDER_CREATED {
		return errors.OrderStatusNotValid
	}

	switch dto.Status {
	case order.ORDER_DELIVERY:
		if orderDAO.Status != order.ORDER_PREPARED {
			return errors.OrderStatusNotValid
		}
		orderDAO.Status = dto.Status

		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			"Đơn hàng được đơn vị vận chuyển tiếp nhận và giao hàng"); err != nil {
			return err
		}
		bodyContent = fmt.Sprintf("Xin chào, đơn hàng [%v] của bạn đã được vận chuyển bởi %v ", orderDAO.OrderID, orderDAO.Delivery.DeliveryName)

	case order.ORDER_CANCEL_BY_USER:
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			fmt.Sprintf("Đơn hàng không thể giao thành công lý do: %v", dto.Message)); err != nil {
			return err
		}
		bodyContent = fmt.Sprintf("Xin chào, đơn hàng [%v] của bạn đã bị huỷ do %v ", orderDAO.OrderID, dto.Message)
	case order.ORDER_CANCEL_BY_DELI:
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			fmt.Sprintf("Đơn hàng giao thất bại lý do: %v", dto.Message)); err != nil {
			return err
		}
		bodyContent = fmt.Sprintf("Xin chào, đơn hàng [%v] của bạn đã bị huỷ do %v ", orderDAO.OrderID, dto.Message)

	case order.ORDER_SHIPPING_FINISH:
		if orderDAO.Status != order.ORDER_DELIVERY {
			return errors.OrderStatusNotValid
		}
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status,
			fmt.Sprintf("Đơn hàng giao thành công: %v", dto.Message)); err != nil {
			return err
		}
		bodyContent = fmt.Sprintf("Xin chào, đơn hàng [%v] của bạn đã giao thành công.", orderDAO.OrderID)

		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_COMPLETED, "Hoàn tất đơn hàng"); err != nil {
			return err
		}

	default:
		return errors.OrderCannotUpdate
	}

	userNoti := msgDTO.NewNotificationMessage(msgDTO.NOTIFY_USER, orderDAO.UserId, "[Latipe] Thông báo tình trạng đơn hàng", bodyContent, orderDAO.Thumbnail)
	//send message to user
	if err := o.notifyPub.NotifyToUser(userNoti); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindByIdForUpdate(ctx, dto.OrderID)
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
	var bodyContent string

	switch dto.Status {
	case msgDTO.ORDER_EVENT_COMMIT_SUCCESS:
		orderDAO, err := o.orderRepo.FindByIdForUpdate(ctx, dto.OrderID)
		if err != nil {
			return err
		}

		bodyContent = "Đơn hàng được tạo thành công"
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_CREATED,
			bodyContent); err != nil {
			return err
		}
		req := dto2.GetStoreByIdRequest{
			StoreID: orderDAO.StoreId,
		}

		storeCms, err := o.storeServ.GetStoreByStoreId(ctx, &req)
		if err != nil {
			return err
		}

		commission := order.OrderCommission{
			OrderType:         "orders",
			OrderID:           orderDAO.OrderID,
			Order:             orderDAO,
			StoreID:           orderDAO.StoreId,
			DiscountFromStore: orderDAO.StoreDiscount,
			Status:            order.COMMIS_PENDING,
			SystemFee:         int(float64(orderDAO.SubTotal) * storeCms.FeePerOrder),
			CreatedAt:         time.Time{},
		}
		commission.AmountReceived = orderDAO.SubTotal - commission.DiscountFromStore - commission.SystemFee

		if err := o.commissionRepo.CreateOrderCommission(&commission); err != nil {
			return err
		}

	case msgDTO.ORDER_EVENT_FAIL_BY_PRODUCT:
		bodyContent = "Đơn hàng xử lý thất bại do lỗi sản phẩm"
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			bodyContent); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PROMOTION:
		bodyContent = "Đơn hàng xử lý thất bại do lỗi khuyến mãi"
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			bodyContent); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_DELIVERY:
		bodyContent = "Đơn hàng xử lý thất bại do lỗi vận chuyển"
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			bodyContent); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PAYMENT:
		bodyContent = "Đơn hàng xử lý thất bại do lỗi thanh toán"
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			bodyContent); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_CANCEL:
		bodyContent = "Đơn hàng bị hủy"
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			bodyContent); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_REFUND:
		bodyContent = "Đơn hàng hoàn trả"
		if err := o.orderRepo.UpdateStatus(ctx, dto.OrderID, order.ORDER_FAILED,
			bodyContent); err != nil {
			return err
		}
	}

	orderDAO, err := o.orderRepo.FindByIdSingleObject(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	userNoti := msgDTO.NewNotificationMessage(msgDTO.NOTIFY_USER, orderDAO.UserId,
		"[Latipe] Thông báo tình trạng đơn hàng", bodyContent, orderDAO.Thumbnail)
	//send message to user
	if err := o.notifyPub.NotifyToUser(userNoti); err != nil {
		return err
	}

	if dto.Status == msgDTO.ORDER_EVENT_COMMIT_SUCCESS {
		bodyContent = fmt.Sprintf("Cửa hàng của bạn vừa có 1 đơn hàng chờ xác nhận từ hệ thống [%v]", orderDAO.OrderID)
		storeNoti := msgDTO.NewNotificationMessage(msgDTO.NOTIFY_STORE, orderDAO.StoreId,
			"[Latipe] Thông báo tình trạng đơn hàng", bodyContent, orderDAO.Thumbnail)

		//send message to user
		if err := o.notifyPub.NotifyToUser(storeNoti); err != nil {
			return err
		}
	}

	return nil
}

func (o orderCommandService) UserRefundOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByIdForUpdate(ctx, dto.OrderID)
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

	mess := msgDTO.OrderCancelMessage{
		OrderID:      dao.OrderID,
		CancelStatus: order.ORDER_REFUND,
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

	mess := msgDTO.OrderCancelMessage{
		OrderID:      dao.OrderID,
		CancelStatus: order.ORDER_CANCEL_BY_USER,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByIdForUpdate(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_CANCEL_BY_ADMIN,
		fmt.Sprintf("Đơn hàng bị hủy do yêu cầu của quản trị viên:%v", dto.Message)); err != nil {
		return err
	}

	mess := msgDTO.OrderCancelMessage{
		OrderID:      dao.OrderID,
		CancelStatus: order.ORDER_CANCEL_BY_ADMIN,
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

func (o orderCommandService) UpdateRatingItem(ctx context.Context, data *msgDTO.RatingMessage) error {
	err := o.orderRepo.UpdateOrderRating(ctx, data.OrderItemId, data.RatingId)
	if err != nil {
		return err
	}

	return nil
}
