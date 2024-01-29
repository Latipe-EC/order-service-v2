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
	"latipe-order-service-v2/pkg/util/mapper"
	"sync"
	"time"

	publishMsg "latipe-order-service-v2/internal/msgqueue"
	"latipe-order-service-v2/pkg/cache/redis"
	"strings"
)

type orderCommandService struct {
	cfg         *config.Config
	orderRepo   order.Repository
	cacheEngine *redis.CacheEngine
	publisher   *publishMsg.PublisherTransactionMessage
	//grpc client
	voucherGrpc  vouchergrpc.VoucherServiceGRPCClient
	productGrpc  productgrpc.ProductServiceGRPCClient
	deliveryGrpc deliverygrpc.DeliveryServiceGRPCClient
	userGrpc     usergrpc.UserServiceGRPCClient
}

func NewOrderCommmandService(cfg *config.Config,
	orderRepo order.Repository,
	cacheEngine *redis.CacheEngine,
	publisher *publishMsg.PublisherTransactionMessage,
	voucherGrpc vouchergrpc.VoucherServiceGRPCClient,
	productGrpc productgrpc.ProductServiceGRPCClient,
	deliveryGrpc deliverygrpc.DeliveryServiceGRPCClient,
	userGrpc usergrpc.UserServiceGRPCClient) OrderCommandUsecase {
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

func (o orderCommandService) genOrderKey(userId string) string {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")[:10]
	key := fmt.Sprintf("%v%v%v", o.cfg.Server.KeyID, userId[:4], keyGen)

	return key
}

func (o orderCommandService) CreateOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error) {

	userReq := usergrpc.GetDetailAddressRequest{AddressId: dto.Address.AddressId}
	address, err := o.userGrpc.GetAddressDetail(ctx, &userReq)
	if err != nil {
		return nil, err
	}

	//get product data from product service
	var orders []*order.Order

	// Handle order group by store
	var orderIDs []string

	for _, i := range dto.StoreOrders {

		productReq := productgrpc.GetPurchaseProductRequest{}
		if err := mapper.BindingStructGrpc(i, &productReq); err != nil {
			return nil, err
		}

		products, err := o.productGrpc.CheckInStock(ctx, &productReq)
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
		orderData, err := o.saveOrderIntoDatabase(ctx, dto, &i, address, shippingDetail, products)
		if err != nil {
			log.Error(err)
			continue
		}

		//add key into response and send order data to orchestration service
		orderIDs = append(orderIDs, orderData.OrderID)
		orders = append(orders, orderData)
	}

	var wg sync.WaitGroup
	errMsgChan := make(chan error, len(orders))

	for index, i := range orders {
		wg.Add(1)
		go func(orderData *order.Order, cartIds []string) {
			defer wg.Done()

			if err := o.publisher.SendOrderCreatedMessage(MappingDataToMessage(orderData, cartIds)); err != nil {
				errMsgChan <- fmt.Errorf("publish order [%v] failed: %v", orderData.OrderID, err)
			}
		}(i, dto.StoreOrders[index].CartIds)
	}

	go func() {
		wg.Wait()
		close(errMsgChan)
	}()

	// Check for errors in the goroutines
	for err := range errMsgChan {
		log.Error(err)
		// handle errors here, such as returning an error response or logging them.
	}

	data := orderDTO.CreateOrderResponse{
		UserOrder: orderDTO.UserRequest{
			UserId:   dto.UserRequest.UserId,
			Username: dto.UserRequest.Username,
		},
		OrderKeys: orderIDs,
		CreatedAt: time.Now(),
	}

	return &data, nil
}

func (o orderCommandService) saveOrderIntoDatabase(ctx context.Context, dto *orderDTO.CreateOrderRequest,
	storeOrder *orderDTO.StoreOrder, address *usergrpc.GetDetailAddressResponse,
	deli *deliverygrpc.GetShippingCostResponse, productItems *productgrpc.GetPurchaseItemResponse) (*order.Order, error) {

	orderDAO := order.Order{}
	orderDAO.Username = dto.UserRequest.Username
	orderDAO.UserId = dto.UserRequest.UserId

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
			StoreID:     item.StoreId,
			NameOption:  item.NameOption,
			OptionID:    item.OptionId,
			Quantity:    int(item.Quantity),
			Price:       int(item.Price),
			NetPrice:    int(item.PromotionalPrice),
			ProdImg:     item.Image,
		}
		//calculate subtotal of item
		if i.NetPrice != 0 {
			i.SubTotal = i.NetPrice * i.Quantity
		} else {
			i.SubTotal = i.Price * i.Quantity
		}
		orderItems = append(orderItems, &i)
	}
	orderDAO.OrderItem = orderItems

	// Check if the store order has a voucher to apply a discount to the order
	if len(storeOrder.VoucherCode) > 0 {

		voucherDetail, err := o.voucherGrpc.CheckingVoucher(ctx,
			orderquery.MappingVoucherRequest(storeOrder.VoucherCode, &orderDAO))
		if err != nil {
			return nil, err
		}

		if voucherDetail.IsSuccess == true {

			for _, v := range voucherDetail.Items {
				switch v.VoucherType {
				case voucherConst.FREE_SHIP:
					if deli.Cost < v.DiscountValue {
						orderDAO.ShippingDiscount = int(deli.Cost)
					} else {
						orderDAO.ShippingDiscount = int(v.DiscountValue)
					}

				case voucherConst.DISCOUNT_ORDER:
					orderDAO.ItemDiscount = int(v.DiscountValue)
				}

			}
		}
		orderDAO.Vouchers = strings.Join(storeOrder.VoucherCode, ";")
	}

	//calculate amount order
	orderDAO.SubTotal += orderDAO.ShippingCost
	orderDAO.Amount = orderDAO.SubTotal - (orderDAO.ItemDiscount + orderDAO.ShippingDiscount)
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

	orderDAO.PaymentMethod = storeOrder.PaymentMethod
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

func (o orderCommandService) CancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if dao.UserId != dto.UserId {
		return errors.ErrNotFoundRecord
	}

	if dao.Status == order.ORDER_CANCEL {
		return errors.ErrNotChange
	}

	if dao.Status != order.ORDER_CREATED {
		return errors.OrderCannotCancel
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_CANCEL); err != nil {
		return err
	}

	mess := msgDTO.OrderStatusMessage{
		OrderID: dao.OrderID,
		Status:  order.ORDER_CANCEL,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
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

	if dao.Status == order.ORDER_REFUND {
		return errors.ErrNotChange
	}

	if dao.Status != order.ORDER_SHIPPING_FINISH {
		return errors.OrderCannotCancel
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_REFUND); err != nil {
		return err
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

func (o orderCommandService) AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if dao.Status == order.ORDER_CANCEL {
		return errors.ErrNotChange
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.OrderID, order.ORDER_CANCEL); err != nil {
		return err
	}

	mess := msgDTO.OrderStatusMessage{
		OrderID: dao.OrderID,
		Status:  order.ORDER_CANCEL,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	if orderDAO.Status == order.ORDER_CANCEL {
		return errors.ErrBadRequest
	}

	orderDAO.Status = dto.Status

	if err := o.orderRepo.Update(ctx, *orderDAO); err != nil {
		return err
	}

	return nil
}

func (o orderCommandService) UpdateOrderStatusByEvent(ctx context.Context, dto *msgDTO.OrderStatusMessage) error {

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return err
	}

	switch dto.Status {
	case msgDTO.ORDER_EVENT_COMMIT_SUCCESS:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_CREATED,
			"Đơn hàng được tạo thành công"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PRODUCT:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi sản phẩm"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PROMOTION:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi khuyến mãi"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_DELIVERY:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi vận chuyển"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_FAIL_BY_PAYMENT:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_FAILED,
			"Đơn hàng xử lý thất bại do lỗi thanh toán"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_CANCEL:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_FAILED,
			"Đơn hàng bị hủy"); err != nil {
			return err
		}
	case msgDTO.ORDER_EVENT_REFUND:
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_FAILED,
			"Đơn hàng hoàn trả"); err != nil {
			return err
		}
	}

	return nil
}

func (o orderCommandService) DeliveryUpdateStatusOrder(ctx context.Context, dto delivery.UpdateOrderStatusRequest) (*delivery.UpdateOrderStatusResponse, error) {
	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	if orderDAO.Status != order.ORDER_DELIVERY {
		return nil, errors.OrderStatusNotValid
	}

	if orderDAO.Delivery.DeliveryId != dto.DeliveryID {
		return nil, errors.ErrNotFoundRecord
	}

	if dto.Status == order.ORDER_CANCEL || dto.Status == order.ORDER_SHIPPING_FINISH {
		orderDAO.Status = dto.Status
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, dto.Status); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.ErrBadRequest
	}

	ordMsg := msgDTO.OrderStatusMessage{
		Status:  dto.Status,
		OrderID: orderDAO.OrderID,
	}

	if dto.Status == order.ORDER_CANCEL || dto.Status == order.ORDER_SHIPPING_FINISH {
		if err := o.publisher.SendOrderCancelMessage(&ordMsg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (o orderCommandService) UpdateOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error) {
	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	notFound := true
	itemPreparedCount := 0

	for _, i := range orderDAO.OrderItem {

		if i.StoreID != dto.StoreId {
			continue
		}

		if i.Id == dto.ItemID {
			notFound = false
			if i.Status != order.OI_PREPARED && i.Status != order.OI_CANCEL {
				if err := o.orderRepo.UpdateOrderItem(ctx, i.Id, order.OI_PREPARED); err != nil {
					return nil, err
				}
				i.Status = order.OI_PREPARED
			} else {
				return nil, errors.ErrNotChange
			}
		}

		if i.Status == order.OI_PREPARED {
			itemPreparedCount++
		}
	}

	if notFound {
		return nil, errors.ErrNotFoundRecord
	}

	if orderDAO.Status == order.ORDER_CREATED {
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_PENDING); err != nil {
			return nil, err
		}
	}

	if len(orderDAO.OrderItem) == itemPreparedCount {
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_DELIVERY); err != nil {
			return nil, err
		}
	}

	resp := store.UpdateOrderItemResponse{
		OrderID: dto.OrderID,
		ItemID:  dto.ItemID,
		Status:  order.OI_PREPARED,
	}

	return &resp, nil
}

func (o orderCommandService) CancelOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error) {
	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	notFound := true

	for _, i := range orderDAO.OrderItem {

		if i.StoreID != dto.StoreId {
			continue
		}

		if i.Status != order.OI_PREPARED && i.Id == dto.ItemID {
			notFound = false
			if err := o.orderRepo.UpdateOrderItem(ctx, i.Id, order.ORDER_CANCEL); err != nil {
				return nil, err
			}
			i.Status = order.OI_PREPARED
		}

	}

	if notFound {
		return nil, errors.ErrNotFoundRecord
	}

	if err := o.orderRepo.UpdateStatus(ctx, orderDAO.OrderID, order.ORDER_CANCEL,
		"Đơn hàng bị hủy do nhà cung cấp không thể chuẩn bị sản phẩm"); err != nil {
		return nil, err
	}

	resp := store.UpdateOrderItemResponse{
		OrderID: dto.OrderID,
		ItemID:  dto.ItemID,
		Status:  order.ORDER_CANCEL,
	}

	return &resp, nil
}

func (o orderCommandService) UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}
