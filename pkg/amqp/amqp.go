package amqp

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Amqp struct {
	Connect *amqp091.Connection `json:"connect"`
}

func New(url string) (*Amqp, error) {
	connAttempts := 10
	var err error
	for connAttempts > 0 {
		conn, err := amqp091.Dial(url)
		if err == nil {
			return &Amqp{Connect: conn}, nil
		}
		slog.Info("Amqp is trying to connect, attempts left: %d", connAttempts)
		time.Sleep(time.Second)
		connAttempts--
	}
	return nil, fmt.Errorf("amqp - New - amqp091.Dial: %w", err)
}
// 	conn, err := amqp091.Dial(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Amqp{Connect: conn}, nil
// }

func (a *Amqp) Close() error {
	return a.Connect.Close()
}