package grpc_adapt

import (
	"github.com/google/wire"
	deliverygrpc "latipe-order-service-v2/internal/infrastructure/grpc/deliveryServ"
	productgrpc "latipe-order-service-v2/internal/infrastructure/grpc/productServ"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
)

var Set = wire.NewSet(
	vouchergrpc.NewVoucherClientGrpcConnection,
	productgrpc.NewProductGrpcClientConnection,
	deliverygrpc.NewDeliveryClientGrpcConnection,
)
