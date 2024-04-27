package amqp

import "github.com/rabbitmq/amqp091-go"

type Amqp struct {
	Connect *amqp091.Connection `json:"connect"`
}

func New(url string) (*Amqp, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Amqp{Connect: conn}, nil
}

func (a *Amqp) Close() error {
	return a.Connect.Close()
}