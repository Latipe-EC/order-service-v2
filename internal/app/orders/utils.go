package orders

import (
	"latipe-order-service-v2/internal/domain/dto/order"
	enitites "latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msg"
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv/dto"
	dto2 "latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
)

func CheckStoreHaveOrder(entities enitites.Order, storeId string) bool {
	for _, i := range entities.OrderItem {
		if i.StoreID == storeId {
			return true
		}
	}
	return false
}

func MappingOrderItemToGetInfo(request *order.CreateOrderRequest) []dto.StoreOrderRequest {
	var storeOrder []dto.StoreOrderRequest
	for _, i := range request.StoreOrders {

		var items []dto.ValidateItems
		for _, j := range i.Items {
			product := dto.ValidateItems{
				ProductId: j.ProductId,
				OptionId:  j.OptionId,
				Quantity:  j.Quantity,
			}
			items = append(items, product)
		}

		ownerProduct := dto.StoreOrderRequest{
			StoreID: i.StoreID,
			Items:   items,
		}

		storeOrder = append(storeOrder, ownerProduct)
	}
	return storeOrder
}

func MappingVoucherRequest(dto *order.CreateOrderRequest, voucherCode []string, orderData *msg.OrderMessage) dto2.CheckingVoucherRequest {
	voucherReq := dto2.CheckingVoucherRequest{}
	voucherReq.Vouchers = voucherCode
	voucherReq.AuthorizationHeader.BearerToken = dto.Header.BearerToken
	voucherReq.OrderTotalAmount = orderData.SubTotal - orderData.ShippingCost
	voucherReq.PaymentMethod = orderData.PaymentMethod
	voucherReq.UserId = orderData.UserRequest.UserId
	return voucherReq
}
