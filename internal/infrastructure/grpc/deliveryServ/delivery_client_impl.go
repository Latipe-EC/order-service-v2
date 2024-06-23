package deliverygrpc

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"latipe-order-service-v2/config"
)

type deliveryServiceGRPCClientImpl struct {
	cfg *config.Config
	cc  grpc.ClientConnInterface
}

func NewDeliveryServiceGRPCClientImpl(config *config.Config) DeliveryServiceClient {
	// Set up a connection to the server.
	log.Info("[GRPC Client] open connection to delivery service")
	conn, err := grpc.Dial(config.GRPC.DeliveryServiceGrpc.Connection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: %v", err)
	}

	return &deliveryServiceGRPCClientImpl{
		cfg: config,
		cc:  conn,
	}
}

func (d deliveryServiceGRPCClientImpl) CalculateShippingCost(ctx context.Context, in *GetShippingCostRequest, opts ...grpc.CallOption) (*GetShippingCostResponse, error) {
	md := metadata.New(
		map[string]string{"x-api-key": d.cfg.GRPC.DeliveryServiceGrpc.APIKey},
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	out := new(GetShippingCostResponse)
	err := d.cc.Invoke(ctx, "/DeliveryService/CalculateShippingCost", in, out, opts...)
	if err != nil {
		log.Errorf("request to gRPC is failed cause %v", err)
		return nil, err
	}

	return out, nil
}
