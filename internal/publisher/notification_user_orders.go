package publisher

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/domain/msgDTO"
	"time"
)

type NotificationMessagePublisher struct {
	channel *amqp.Channel
	cfg     *config.Config
	conn    *amqp.Connection
}

func NewNotificationMessagePublisher(cfg *config.Config, conn *amqp.Connection) *NotificationMessagePublisher {
	producer := NotificationMessagePublisher{
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

func (pub *NotificationMessagePublisher) NotifyToUser(notify *msgDTO.SendNotificationMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToByte(&notify)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		pub.cfg.RabbitMQ.NotificationEvent.Exchange,
		pub.cfg.RabbitMQ.NotificationEvent.RoutingKey)

	err = pub.channel.PublishWithContext(ctx,
		pub.cfg.RabbitMQ.NotificationEvent.Exchange,   // exchange
		pub.cfg.RabbitMQ.NotificationEvent.RoutingKey, // routing key
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
