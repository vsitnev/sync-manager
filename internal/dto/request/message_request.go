package dto

import "github.com/vsitnev/sync-manager/internal/model"

type MessageRequestDto struct {
	Routing string            `json:"routing" db:"routing"`
	Message model.AmqpMessage `json:"message" db:"message"`
}

func (m *MessageRequestDto) ToModel() model.Message {
	return model.Message{
		Routing: m.Routing,
		Message: m.Message,
	}
}