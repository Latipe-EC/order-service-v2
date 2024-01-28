package rabbitclient

import (
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-order-service-v2/config"
)

func NewRabbitClientConnection(globalCfg *config.Config) *amqp.Connection {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": globalCfg.RabbitMQ.NameService,
		},
	}

	conn, err := amqp.DialConfig(globalCfg.RabbitMQ.Connection, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ cause:%v", err)
	}

	log.Info("Comsumer has been connected")
	return conn
}
