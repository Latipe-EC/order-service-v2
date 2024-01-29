package usergrpc

import (
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"latipe-order-service-v2/config"
)

func NewUserServiceClientConnection(config *config.Config) UserServiceGRPCClient {
	log.Info("[GRPC Client] open connection to user service")
	conn, err := grpc.Dial(config.GRPC.UserConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}

	return NewUserServiceGRPCClient(conn)
}
