package orders

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"latipe-order-service-v2/internal/common/errors"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msg"
	deliDto "latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv/dto"
	prodServDTO "latipe-order-service-v2/internal/infrastructure/adapter/productserv/dto"
	userDTO "latipe-order-service-v2/internal/infrastructure/adapter/userserv/dto"
	voucherDTO "latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
	"latipe-order-service-v2/pkg/util/mapper"
	"strings"
	"time"
)

func (o orderService) CreateOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error) {

	//get product data from product service
	productReq := prodServDTO.OrderProductRequest{
		StoreOrders: MappingOrderItemToGetInfo(dto),
	}
	products, err := o.productServ.GetProductOrderInfo(ctx, &productReq)
	if err != nil {
		return nil, err
	}

	//get address data from user service
	addressRequest := userDTO.GetDetailAddressRequest{
		AddressId:           dto.Address.AddressId,
		AuthorizationHeader: userDTO.AuthorizationHeader{BearerToken: dto.Header.BearerToken},
	}
	userAddress, err := o.userServ.GetAddressDetails(ctx, &addressRequest)
	if err != nil {
		return nil, err
	}

	var orders []*order.Order

	// Handle order group by store
	var orderIDs []string
	for _, i := range products.Items {

		// Find the store order (user request) that matches the product service request
		var storeOrder orderDTO.StoreOrder

		for index, j := range dto.StoreOrders {
			if j.StoreID == i.StoreID {
				storeOrder = j
				// Remove the found store order from the slice
				dto.StoreOrders = append(dto.StoreOrders[:index], dto.StoreOrders[index+1:]...)
				break
			}
			continue
		}

		//handle order group by store_order and calculate shipping cost
		shippingReq := deliDto.GetShippingCostRequest{
			SrcCode:    i.ProvinceCode,
			DestCode:   userAddress.CityOrProvinceId,
			DeliveryId: storeOrder.Delivery.DeliveryId,
		}
		shippingDetail, err := o.deliServ.CalculateShippingCost(ctx, &shippingReq)
		if err != nil {
			log.Error(err)
			continue
		}

		//init order data
		orderData, err := o.saveOrderIntoDatabase(ctx, dto, userAddress, shippingDetail, &i, &storeOrder)
		if err != nil {
			log.Error(err)
			continue
		}

		//add key into response and send order data to orchestration service
		orderIDs = append(orderIDs, orderData.OrderID)
		orders = append(orders, orderData)
	}

	if err := o.publisher.SendOrderCreatedMessage(orders); err != nil {
		return nil, err
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

func (o orderService) saveOrderIntoDatabase(ctx context.Context, dto *orderDTO.CreateOrderRequest,
	address *userDTO.GetDetailAddressResponse, deli *deliDto.GetShippingCostResponse,
	productDetails *prodServDTO.StoreOrder, storeOrder *orderDTO.StoreOrder) (*order.Order, error) {

	orderDAO := order.Order{}
	orderDAO.Username = dto.UserRequest.Username
	orderDAO.UserId = dto.UserRequest.UserId

	if err := mapper.BindingStruct(dto, &orderDAO); err != nil {
		log.Errorf("Mapping value from dto to dao failed cause: %s", err)
		return nil, err
	}
	orderDAO.ShippingCost = deli.Cost

	var orderItems []*order.OrderItem
	for _, item := range productDetails.Items {
		i := order.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.Name,
			StoreID:     item.StoreId,
			NameOption:  item.NameOption,
			OptionID:    item.OptionId,
			Quantity:    item.Quantity,
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

		voucherReq := MappingVoucherRequest(dto, storeOrder.VoucherCode, &orderDAO)
		voucherDetail, err := o.voucherSer.CheckingVoucher(ctx, &voucherReq)
		if err != nil {
			return nil, err
		}

		if voucherDetail.IsSuccess == true {

			for _, v := range voucherDetail.Items {
				switch v.VoucherType {
				case voucherDTO.FREE_SHIP:
					if deli.Cost < v.DiscountValue {
						orderDAO.ShippingDiscount = deli.Cost
					} else {
						orderDAO.ShippingDiscount = v.DiscountValue
					}

				case voucherDTO.DISCOUNT_ORDER:
					orderDAO.ItemDiscount = v.DiscountValue
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
		Cost:            deli.Cost,
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

func (o orderService) CancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
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

	mess := msg.OrderMessage{
		OrderID:       dao.OrderID,
		Status:        order.ORDER_CANCEL,
		PaymentMethod: dao.PaymentMethod,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderService) UserRefundOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
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

	mess := msg.OrderMessage{
		OrderID:       dao.OrderID,
		Status:        order.ORDER_REFUND,
		PaymentMethod: dao.PaymentMethod,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderService) AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
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

	mess := msg.OrderMessage{
		OrderID:       dao.OrderID,
		Status:        order.ORDER_CANCEL,
		PaymentMethod: dao.PaymentMethod,
	}

	if err := o.publisher.SendOrderCancelMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

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

func (o orderService) DeliveryUpdateStatusOrder(ctx context.Context, dto delivery.UpdateOrderStatusRequest) (*delivery.UpdateOrderStatusResponse, error) {
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

	ordMsg := msg.OrderMessage{
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

func (o orderService) UpdateOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error) {
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

func (o orderService) CancelOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error) {
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

func (o orderService) UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}
