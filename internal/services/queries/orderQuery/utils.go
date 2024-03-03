package orderQuery

import (
	"golang.org/x/exp/slices"
	dto "latipe-order-service-v2/internal/domain/dto/order"
	enitites "latipe-order-service-v2/internal/domain/entities/order"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
)

func MappingShippingVoucherRequest(req *dto.PromotionData, orderData *enitites.Order, storeId string) *vouchergrpc.CheckoutVoucherRequest {
	if slices.ContainsFunc(req.FreeShippingVoucherInfo.StoreIds, func(s string) bool {
		return s == storeId
	}) {
		voucherReq := vouchergrpc.CheckoutVoucherRequest{
			OrderTotalAmount: int64(orderData.SubTotal),
			PaymentMethod:    int32(orderData.PaymentMethod),
			UserId:           orderData.UserId,
		}

		freeshipVoucher := vouchergrpc.VoucherData{
			VoucherCode: req.FreeShippingVoucherInfo.VoucherCode,
		}
		voucherReq.VoucherData = &freeshipVoucher
		return &voucherReq
	}

	return nil
}

func MappingPaymentVoucherRequest(req *dto.PromotionData, subTotal int64, paymentMethod int, userId string) *vouchergrpc.CheckoutVoucherRequest {

	voucherReq := vouchergrpc.CheckoutVoucherRequest{
		OrderTotalAmount: subTotal,
		PaymentMethod:    int32(paymentMethod),
		UserId:           userId,
	}

	paymentVoucher := vouchergrpc.VoucherData{
		VoucherCode: req.PaymentVoucherInfo.VoucherCode,
	}
	voucherReq.VoucherData = &paymentVoucher

	return &voucherReq
}

func MappingShopVoucherRequest(orderData *enitites.Order, voucherCode string) *vouchergrpc.CheckoutVoucherRequest {
	shopVoucher := vouchergrpc.VoucherData{
		VoucherCode: voucherCode,
	}

	voucherReq := vouchergrpc.CheckoutVoucherRequest{
		OrderTotalAmount: int64(orderData.SubTotal),
		PaymentMethod:    int32(orderData.PaymentMethod),
		UserId:           orderData.UserId,
		VoucherData:      &shopVoucher,
	}

	return &voucherReq
}
