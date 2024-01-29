package deliverygrpc

import (
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"latipe-order-service-v2/config"
)

func NewDeliveryClientGrpcConnection(config *config.Config) DeliveryServiceGRPCClient {
	// Set up a connection to the server.
	log.Info("[GRPC Client] open connection to delivery service")
	conn, err := grpc.Dial(config.GRPC.DeliveryConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect: %v", err)
	}

	return NewDeliveryServiceGRPCClient(conn)
}
