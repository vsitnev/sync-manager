package service

import (
	"context"

	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/pkg/amqp"
)

type MessageInput struct {
	Source   string
	Routing  string
	SortType string
	Offset   int
	Limit    int
}

type Message interface {
	SendMessage(ctx context.Context, message model.Message) (int, error)
	GetMessages(ctx context.Context, input MessageInput) ([]model.Message, error)
	GetMessageByID(ctx context.Context, ID int) (model.Message, error)
}

type ServiceDeps struct {
	Reps *repository.Repositories
	Amqp *amqp.Amqp
}

type Services struct {
	Message Message
}

func NewServices(deps ServiceDeps) *Services {
	return &Services{
		Message: NewMessageService(deps.Reps, deps.Amqp),
	}
}
