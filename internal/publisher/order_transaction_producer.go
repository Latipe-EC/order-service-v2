package publisher

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/domain/msgDTO"
	"time"
)

type PublisherTransactionMessage struct {
	channel *amqp.Channel
	cfg     *config.Config
	conn    *amqp.Connection
}

func NewTransactionProducer(cfg *config.Config, conn *amqp.Connection) *PublisherTransactionMessage {
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

func (pub *PublisherTransactionMessage) SendOrderCreatedMessage(orderMsg *msgDTO.CreateOrderMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&orderMsg)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEvent.PublishRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderEvent.Exchange,          // exchange
		pub.cfg.RabbitMQ.SagaOrderEvent.PublishRoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Errorf("Publishing message was failed cause:%v", err)
		return err
	}

	return nil
}

func (pub *PublisherTransactionMessage) SendOrderCancelMessage(request *msgDTO.OrderCancelMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&request)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.SagaOrderEvent.Exchange,
		pub.cfg.RabbitMQ.SagaOrderEvent.CancelRoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.SagaOrderEvent.Exchange,         // exchange
		pub.cfg.RabbitMQ.SagaOrderEvent.CancelRoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	return nil
}
