package purchase

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/domain/msgDTO"
	"latipe-order-service-v2/internal/services/commands/orderCmd"
	"sync"
	"time"
)

type PurchaseReplySubscriber struct {
	config    *config.Config
	conn      *amqp.Connection
	orderServ orderCmd.OrderCommandUsecase
}

func NewPurchaseReplySubscriber(cfg *config.Config, conn *amqp.Connection, orderServ orderCmd.OrderCommandUsecase) *PurchaseReplySubscriber {
	return &PurchaseReplySubscriber{
		config:    cfg,
		orderServ: orderServ,
		conn:      conn,
	}
}
func (mq PurchaseReplySubscriber) ListenOrderEventQueue(wg *sync.WaitGroup) {
	channel, err := mq.conn.Channel()
	defer channel.Close()

	// define an exchange type "topic"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare exchange: %v", err)
	}

	// create queue
	q, err := channel.QueueDeclare(
		"transaction_reply",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare queue: %v", err)
	}

	err = channel.QueueBind(
		q.Name,
		mq.config.RabbitMQ.SagaOrderEvent.ReplyRoutingKey,
		mq.config.RabbitMQ.SagaOrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                         // queue
		mq.config.RabbitMQ.NameService, // consumer
		true,                           // auto ack
		false,                          // exclusive
		false,                          // no local
		false,                          // no wait
		nil,                            //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := mq.messageHandler(msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}

	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")
}

func (mq PurchaseReplySubscriber) messageHandler(msg amqp.Delivery) error {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := msgDTO.OrderStatusMessage{}

	if err := json.Unmarshal(msg.Body, &message); err != nil {
		log.Infof("Parse message to order failed cause: %s", err)
		return err
	}

	if err := mq.orderServ.UpdateOrderStatusByReplyMessage(ctx, &message); err != nil {
		log.Infof("Handling order message was failed cause: %s", err)
		return err
	}

	endTime := time.Now()
	log.Infof("The order [%v]  was processed successfully - duration:%v", message.OrderID, endTime.Sub(startTime))
	return nil
}
