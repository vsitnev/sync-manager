package dto

import (
	"time"
)

type Message struct {
	ID      int         `json:"id" db:"id"`
	Routing string      `json:"routing" db:"routing"`
	Message AmqpMessage `json:"message" db:"message"`
	Dead    bool        `json:"dead" db:"dead"`
	Retried bool        `json:"retried" db:"retried"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type AmqpMessage struct {
	MessageID string      `json:"message_id" db:"id" binding:"required"`
	Source    string      `json:"source" db:"source" binding:"required"`
	Operation string      `json:"operation" db:"operation" binding:"required"`
	Created   int64       `json:"created" db:"created" binding:"required"`
	Data      interface{} `json:"data" db:"data" binding:"required"`
}
