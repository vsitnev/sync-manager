package model

import "time"

type Message struct {
	ID        int         `json:"id" db:"id"`
	Routing   string      `json:"routing" db:"routing"`
	Message   AmqpMessage `json:"message" db:"message"`
	Dead      bool        `json:"dead" db:"dead"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}

// FIXME: Data field must be shown as empty object in swagger schema
type AmqpMessage struct {
	ID        string                 `json:"id" db:"id"`
	Operation string                 `json:"routing" db:"routing"`
	Created   int32                  `json:"message" db:"message"`
	Data      struct{}               `json:"dead,omitempty" db:"dead" swaggertype:"object,{}="`
}
