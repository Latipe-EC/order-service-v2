package vouchergrpc

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"latipe-order-service-v2/config"
)

func NewVoucherClientGrpcImpl(config *config.Config) VoucherServiceClient {
	// Set up a connection to the server.
	log.Info("[GRPC Client] open connection to promotion service")
	conn, err := grpc.Dial(config.GRPC.VoucherServiceGrpc.Connection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	return &voucherServiceGRPCClientImpl{
		cfg: config,
		cc:  conn,
	}
}

type voucherServiceGRPCClientImpl struct {
	cfg *config.Config
	cc  grpc.ClientConnInterface
}

func (v voucherServiceGRPCClientImpl) CheckoutVoucherForPurchase(ctx context.Context, in *CheckoutVoucherRequest, opts ...grpc.CallOption) (*CheckoutVoucherResponse, error) {
	md := metadata.New(
		map[string]string{"x-api-key": v.cfg.GRPC.VoucherServiceGrpc.APIKey},
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	out := new(CheckoutVoucherResponse)
	err := v.cc.Invoke(ctx, "/VoucherService/CheckoutVoucherForPurchase", in, out, opts...)
	if err != nil {
		log.Errorf("request to gRPC is failed cause %v", err)
		return nil, err
	}
	return out, nil
}

func (v voucherServiceGRPCClientImpl) ApplyVoucherToPurchase(ctx context.Context, in *ApplyVoucherRequest, opts ...grpc.CallOption) (*ApplyVoucherResponse, error) {
	md := metadata.New(
		map[string]string{"x-api-key": v.cfg.GRPC.VoucherServiceGrpc.APIKey},
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	out := new(ApplyVoucherResponse)
	err := v.cc.Invoke(ctx, "/VoucherService/ApplyVoucherToPurchase", in, out, opts...)
	if err != nil {
		log.Errorf("request to gRPC is failed cause %v", err)
		return nil, err
	}
	return out, nil
}
