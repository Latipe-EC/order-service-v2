package orderquery

import (
	enitites "latipe-order-service-v2/internal/domain/entities/order"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
)

func MappingVoucherRequest(voucherCode []string, orderData *enitites.Order) *vouchergrpc.CheckingVoucherRequest {
	voucherReq := vouchergrpc.CheckingVoucherRequest{
		Vouchers:         voucherCode,
		OrderTotalAmount: int64(orderData.SubTotal - orderData.ShippingCost),
		PaymentMethod:    int32(orderData.PaymentMethod),
		UserId:           orderData.UserId,
	}
	return &voucherReq
}
