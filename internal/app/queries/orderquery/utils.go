package orderquery

import (
	dto "latipe-order-service-v2/internal/domain/dto/order"
	enitites "latipe-order-service-v2/internal/domain/entities/order"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
)

func MappingPaymentAndShippingVoucherRequest(req *dto.CreateOrderRequest, orderData *enitites.Order) *vouchergrpc.CheckoutVoucherRequest {

	paymentVoucher := vouchergrpc.PaymentVoucher{
		VoucherCode: req.PromotionData.PaymentVoucherInfo.VoucherCode,
	}

	freeshipVoucher := vouchergrpc.FreeShippingVoucher{
		VoucherCode: req.PromotionData.FreeShippingVoucherInfo.VoucherCode,
	}

	voucherReq := vouchergrpc.CheckoutVoucherRequest{
		OrderTotalAmount:    int64(orderData.SubTotal),
		PaymentMethod:       int32(req.PaymentMethod),
		FreeShippingVoucher: &freeshipVoucher,
		PaymentVoucher:      &paymentVoucher,
		UserId:              orderData.UserId,
	}

	return &voucherReq
}

func MappingShopVoucherRequest(orderData *enitites.Order, voucherCode string) *vouchergrpc.CheckoutVoucherRequest {

	voucherReq := vouchergrpc.CheckoutVoucherRequest{
		OrderTotalAmount: int64(orderData.SubTotal),
		PaymentMethod:    int32(orderData.PaymentMethod),
		UserId:           orderData.UserId,
		ShopVouchers: []*vouchergrpc.ShopVoucher{
			{SubTotal: int64(orderData.SubTotal), VoucherCode: voucherCode},
		},
	}

	return &voucherReq
}
