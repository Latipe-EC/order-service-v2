package productgrpc

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"latipe-order-service-v2/config"
)

type productServiceGRPCClientImpl struct {
	cfg *config.Config
	cc  grpc.ClientConnInterface
}

func NewProductGrpcClientImpl(config *config.Config) ProductServiceClient {
	// Set up a connection to the server.
	log.Info("[GRPC Client] open connection to product service")
	conn, err := grpc.Dial(config.GRPC.ProductServiceGrpc.Connection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	return &productServiceGRPCClientImpl{
		cfg: config,
		cc:  conn,
	}
}

func (p productServiceGRPCClientImpl) CheckInStock(ctx context.Context, in *GetPurchaseProductRequest, opts ...grpc.CallOption) (*GetPurchaseItemResponse, error) {
	md := metadata.New(
		map[string]string{"x-api-key": p.cfg.GRPC.ProductServiceGrpc.APIKey},
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	out := new(GetPurchaseItemResponse)
	err := p.cc.Invoke(ctx, "/ProductService/CheckInStock", in, out, opts...)
	if err != nil {
		log.Errorf("request to gRPC is failed cause %v", err)
		return nil, err
	}
	return out, nil
}

func (p productServiceGRPCClientImpl) UpdateQuantity(ctx context.Context, in *UpdateProductQuantityRequest, opts ...grpc.CallOption) (*UpdateProductQuantityResponse, error) {
	out := new(UpdateProductQuantityResponse)
	err := p.cc.Invoke(ctx, "/protobuf.ProductService/UpdateQuantity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
