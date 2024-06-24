package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/pkg/amqpclient"
	"log"
	"sync"
	"time"
)

type MessageConsumer struct {
	amqp    *amqpclient.Amqp
	msgChan <-chan amqp091.Delivery
	handler func(context.Context, []model.Message) error
	wg      sync.WaitGroup
}

func NewMessageConsumer(amqp *amqpclient.Amqp, handler func(context.Context, []model.Message) error) (*MessageConsumer, error) {
	msgChan, err := amqp.Consume("navi.sync")
	if err != nil {
		return nil, fmt.Errorf("MessageConsumer.NewMessageConsumer - amqp.Consume: %v", err)
	}

	c := &MessageConsumer{
		amqp:    amqp,
		msgChan: msgChan,
		handler: handler,
	}

	go c.start()

	return c, nil
}

func (c *MessageConsumer) start() {
	var messages []model.Message

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	/**
	давай продумаем как будем обрабатывать:
	1. navi.sync - navi.sync.[sub-system-name].[domain]
		- запрос будет осуществляться только на Http [POST].
		- маршрутизация исходя из RoutingKey.
	2. POST /messages
		- запрос осуществляется только на AMQP.
		- маршрутизируется исходя из Routing.
	*/
	input := make(chan []model.Message)
	go c.handleMessage(input)

	for {
		select {
		case msg := <-c.msgChan:
			fmt.Println("Message received")
			var amqpMsg model.AmqpMessage
			err := json.Unmarshal(msg.Body, &amqpMsg)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(amqpMsg)
			messages = append(messages, model.Message{
				Message: amqpMsg,
				Routing: msg.RoutingKey,
				Dead:    false,
				Retried: false,
			})

		case <-ticker.C:
			if len(messages) > 0 {
				input <- messages
				messages = nil
			}
		}
	}
}

func (c *MessageConsumer) handleMessage(input chan []model.Message) {
	for msg := range input {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := c.handler(ctx, msg); err != nil {
			log.Println(err)
		}
	}
	close(input)
}

func (c *MessageConsumer) Close() error {
	return c.amqp.Close()
}
