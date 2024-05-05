package service

import (
	"context"

	dto "github.com/vsitnev/sync-manager/internal/dto/request"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/pkg/amqp"
)

type Message interface {
	SendMessage(ctx context.Context, message dto.MessageRequestDto) (int, error)
	GetMessages(ctx context.Context) ([]model.Message, error)
}


type ServiceDeps struct{
	Reps *repository.Repositories
	Amqp *amqp.Amqp
}

type Serices struct {
	Message Message
}

func NewServices(deps ServiceDeps) *Serices {
	return &Serices{
		Message: NewMessageService(deps.Reps, deps.Amqp),
	}
}
