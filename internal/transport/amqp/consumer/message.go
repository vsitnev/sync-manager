package amqp

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/pkg/amqpclient"
	"log"
	"sync"
)

type MessageConsumer struct {
	amqp    *amqpclient.Amqp
	msgChan <-chan amqp091.Delivery
	handler func(context.Context, model.AmqpMessage) error
	wg      sync.WaitGroup
}

func NewMessageConsumer(amqp *amqpclient.Amqp, handler func(context.Context, model.AmqpMessage) error) (*MessageConsumer, error) {
	msgChan, err := amqp.Consume("navi.sync")
	if err != nil {
		return nil, err
	}

	c := &MessageConsumer{
		amqp:    amqp,
		msgChan: msgChan,
		handler: handler,
	}

	c.start()

	return c, nil
}

func (c *MessageConsumer) start() {
	for i := 0; i < c.amqp.PrefetchCount; i++ {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for msg := range c.msgChan {
				log.Printf("Received message: %+v", msg)
				var amqpMsg model.AmqpMessage
				err := json.Unmarshal(msg.Body, &amqpMsg)
				if err != nil {
					log.Println(err)
				}

				if err := c.handler(context.Background(), amqpMsg); err != nil {
					log.Println(err)
				}
				_ = msg.Ack(false)
			}
		}()
	}
}

func (c *MessageConsumer) Close() error {
	c.wg.Wait()
	return c.amqp.Close()
}
