package vouchergrpc

import (
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"latipe-order-service-v2/config"
)

func NewVoucherClientGrpcConnection(config *config.Config) VoucherServiceGRPCClient {
	// Set up a connection to the server.
	log.Info("[GRPC Client] open connection to promotion service")
	conn, err := grpc.Dial(config.GRPC.VoucherConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}

	return NewVoucherServiceGRPCClient(conn)
}
