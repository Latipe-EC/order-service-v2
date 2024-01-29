package productgrpc

import (
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"latipe-order-service-v2/config"
)

func NewProductGrpcClientConnection(config *config.Config) ProductServiceGRPCClient {
	// Set up a connection to the server.
	log.Info("[GRPC Client] open connection to product service")
	conn, err := grpc.Dial(config.GRPC.ProductConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	return NewProductServiceGRPCClient(conn)
}
