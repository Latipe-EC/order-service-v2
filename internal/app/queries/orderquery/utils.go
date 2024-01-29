package orderquery

import (
	"latipe-order-service-v2/internal/domain/dto/order"
	enitites "latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv/dto"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
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

func MappingVoucherRequest(voucherCode []string, orderData *enitites.Order) *vouchergrpc.CheckingVoucherRequest {
	voucherReq := vouchergrpc.CheckingVoucherRequest{
		Vouchers:         voucherCode,
		OrderTotalAmount: int64(orderData.SubTotal - orderData.ShippingCost),
		PaymentMethod:    int32(orderData.PaymentMethod),
		UserId:           orderData.UserId,
	}
	return &voucherReq
}
