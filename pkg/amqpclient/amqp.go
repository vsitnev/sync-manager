package amqpclient

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RoutingKey string
type Queue string
type Exchange struct {
	Name   string
	Routes map[RoutingKey]Queue
}
type Config struct {
	Url           string
	PrefetchCount int
	Exchanges     []Exchange
}

type Amqp struct {
	conn          *amqp091.Connection
	channel       *amqp091.Channel
	PrefetchCount int
}

func New(conf Config) (*Amqp, error) {
	connAttempts := 10
	var err error
	for connAttempts > 0 {
		conn, err := amqp091.Dial(conf.Url)
		if err == nil {
			ch, err := conn.Channel()
			if err != nil {
				return nil, fmt.Errorf("amqpclient - New - conn.Channel: %w", err)
			}
			err = ch.Qos(conf.PrefetchCount, 0, false)
			if err != nil {
				return nil, fmt.Errorf("amqpclient - New - ch.Qos: %w", err)
			}

			err = initQueues(ch, conf.Exchanges)
			if err != nil {
				return nil, fmt.Errorf("amqpclient - New - initQueues: %w", err)
			}

			return &Amqp{conn: conn, channel: ch}, nil
		}
		fmt.Printf("Amqp is trying to connect, attempts left: %d\n", connAttempts)
		time.Sleep(time.Second)
		connAttempts--
	}
	return nil, fmt.Errorf("amqpclient - New - amqp091.Dial: %w", err)
}

func initQueues(ch *amqp091.Channel, exchanges []Exchange) error {
	// For each exchange
	for _, exchange := range exchanges {
		err := ch.ExchangeDeclare(
			exchange.Name, // name
			"direct",      // type
			true,          // durable
			false,         // autoDelete
			false,         // internal
			false,         // noWait
			nil,           // arguments
		)
		if err != nil {
			return fmt.Errorf("amqpclient - New - ch.ExchangeDeclare: %w", err)
		}

		// For each route
		for route, queue := range exchange.Routes {
			_, err = ch.QueueDeclare(
				string(queue), // name
				true,          // durable
				false,         // autoDelete
				false,         // exclusive
				false,         // noWait
				nil,           // arguments
			)
			if err != nil {
				return fmt.Errorf("amqpclient - New - ch.QueueDeclare: %w", err)
			}

			// Bind the queue to the exchange
			err = ch.QueueBind(
				string(queue), // queue name
				string(route), // routing key
				exchange.Name, // exchange name
				false,         // noWait
				nil,           // arguments
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Amqp) Publish(ctx context.Context, exchange, routingKey string, msg []byte) error {
	// Publish the message to the exchange
	//ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	//defer cancel()
	err := a.channel.PublishWithContext(
		ctx,
		exchange,   // exchange name
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
	return err
}

func (a *Amqp) Consume(queue string) (<-chan amqp091.Delivery, error) {
	deliveryChan, err := a.channel.Consume(
		queue, // queue name
		"",    // consumer tag
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return deliveryChan, nil
}

//func (a *Amqp) StartConsumer(queue string, callback func()) {
//	// Start consuming from the queue
//	deliveryChan, err := a.channel.Consume(
//		queue, // queue name
//		"",    // consumer tag
//		true,  // autoAck
//		false, // exclusive
//		false, // noLocal
//		false, // noWait
//		nil,   // arguments
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	go func() {
//		for delivery := range deliveryChan {
//			var message interface{}
//			err := json.Unmarshal(delivery.Body, &message)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//
//			// Process the message here
//			log.Printf("Received message: %+v", message)
//			callback()
//
//			// Acknowledge the message
//			_ = delivery.Ack(false)
//		}
//	}()
//}

func (a *Amqp) Close() error {
	return a.conn.Close()
}
