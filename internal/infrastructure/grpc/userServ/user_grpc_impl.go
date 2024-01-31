package usergrpc

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"latipe-order-service-v2/config"
)

func NewUserServiceClientGRPCImpl(config *config.Config) UserServiceClient {
	log.Info("[GRPC Client] open connection to user service")
	conn, err := grpc.Dial(config.GRPC.UserServiceGrpc.Connection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}

	return &userServiceClientGRPCImpl{
		cfg: config,
		cc:  conn,
	}
}

type userServiceClientGRPCImpl struct {
	cfg *config.Config
	cc  grpc.ClientConnInterface
}

func (u userServiceClientGRPCImpl) GetAddressDetail(ctx context.Context, in *GetDetailAddressRequest, opts ...grpc.CallOption) (*GetDetailAddressResponse, error) {
	md := metadata.New(
		map[string]string{"x-api-key": u.cfg.GRPC.UserServiceGrpc.APIKey},
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	out := new(GetDetailAddressResponse)
	err := u.cc.Invoke(ctx, "/UserService/GetAddressDetail", in, out, opts...)
	if err != nil {
		log.Errorf("request to gRPC is failed cause %v", err)
		return nil, err
	}
	return out, nil
}
