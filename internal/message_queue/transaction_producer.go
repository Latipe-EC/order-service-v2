package message_queue

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-order-service-v2/internal/domain/entities/order"

	"latipe-order-service-v2/config"
	"time"
)

type PublisherTransactionMessage struct {
	channel *amqp.Channel
	cfg     *config.Config
}

func NewTransactionProducer(cfg *config.Config) *PublisherTransactionMessage {
	conn, err := amqp.Dial(cfg.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Info("producer has been connected")

	producer := PublisherTransactionMessage{
		cfg: cfg,
	}

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return nil
	}
	producer.channel = ch

	return &producer
}

func (pub *PublisherTransactionMessage) SendOrderCreatedMessage(data []*order.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, i := range data {
		message := MappingDataToMessage(i)
		body, err := ParseOrderToByte(&message)
		if err != nil {
			return err
		}

		log.Infof("Send message_queue to queue %v - %v",
			pub.cfg.RabbitMQ.OrderCreatedEvent.Exchange,
			pub.cfg.RabbitMQ.OrderCreatedEvent.RoutingKey)

		err = pub.channel.PublishWithContext(ctx,
			pub.cfg.RabbitMQ.OrderCreatedEvent.Exchange,   // exchange
			pub.cfg.RabbitMQ.OrderCreatedEvent.RoutingKey, // routing key
			false,                                         // mandatory
			false,                                         // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		failOnError(err, "Failed to publish a message_queue")
	}

	return nil
}

func (pub *PublisherTransactionMessage) SendOrderCancelMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&request)
	if err != nil {
		return err
	}

	log.Infof("Send message_queue to queue %v - %v",
		pub.cfg.RabbitMQ.OrderCreatedEvent.Exchange,
		pub.cfg.RabbitMQ.OrderCreatedEvent.RoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.OrderCreatedEvent.Exchange,   // exchange
		pub.cfg.RabbitMQ.OrderCreatedEvent.RoutingKey, // routing key
		false,                                         // mandatory
		false,                                         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message_queue")

	return nil
}
