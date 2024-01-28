package msgqueue

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/domain/entities/order"
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

func (pub *PublisherTransactionMessage) SendOrderCreatedMessage(data []*order.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, i := range data {
		message := MappingDataToMessage(i)
		body, err := ParseOrderToByte(&message)
		if err != nil {
			return err
		}

		log.Infof("Send msgqueue to queue %v - %v",
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
		failOnError(err, "Failed to publish a msgqueue")
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

	log.Infof("Send msgqueue to queue %v - %v",
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
	failOnError(err, "Failed to publish a msgqueue")

	return nil
}
