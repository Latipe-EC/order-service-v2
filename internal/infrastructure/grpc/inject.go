package grpc_adapt

import (
	"github.com/google/wire"
	deliverygrpc "latipe-order-service-v2/internal/infrastructure/grpc/deliveryServ"
	productgrpc "latipe-order-service-v2/internal/infrastructure/grpc/productServ"
	vouchergrpc "latipe-order-service-v2/internal/infrastructure/grpc/promotionServ"
	usergrpc "latipe-order-service-v2/internal/infrastructure/grpc/userServ"
)

var Set = wire.NewSet(
	vouchergrpc.NewVoucherClientGrpcImpl,
	productgrpc.NewProductGrpcClientImpl,
	deliverygrpc.NewDeliveryServiceGRPCClientImpl,
	usergrpc.NewUserServiceClientGRPCImpl,
)
