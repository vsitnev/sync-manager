package model

import (
	"database/sql"
	"github.com/vsitnev/sync-manager/internal/dto"
	_ "github.com/vsitnev/sync-manager/pkg/tps"
	"time"
)

type Message struct {
	ID      int         `json:"id" db:"id"`
	Routing string      `json:"routing" db:"routing"`
	Message AmqpMessage `json:"message" db:"message"`
	Dead    bool        `json:"dead" db:"dead"`
	Retried bool        `json:"retried" db:"retried"`

	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}

func (m *Message) ToDto() dto.Message {
	var updated *time.Time
	if m.UpdatedAt.Valid {
		updated = &m.UpdatedAt.Time
	} else {
		updated = nil
	}

	return dto.Message{
		ID:      m.ID,
		Routing: m.Routing,
		Message: dto.AmqpMessage{
			MessageID: m.Message.MessageID,
			Source:    m.Message.Source,
			Operation: m.Message.Operation,
			Created:   m.Message.Created,
			Data:      m.Message.Data,
		},
		Dead:      m.Dead,
		Retried:   m.Retried,
		CreatedAt: m.CreatedAt,
		UpdatedAt: updated,
	}
}

type AmqpMessage struct {
	MessageID string      `json:"message_id" db:"id" binding:"required"`
	Source    string      `json:"source" db:"source" binding:"required"`
	Operation string      `json:"operation" db:"operation" binding:"required"`
	Created   int64       `json:"created" db:"created" binding:"required"`
	Data      interface{} `json:"data" db:"data" binding:"required"`
}
