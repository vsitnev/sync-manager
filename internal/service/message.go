package service

import (
	"context"

	dto "github.com/vsitnev/sync-manager/internal/dto/request"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/pkg/amqp"
)

type MessageService struct {
	repo repository.Message
	amqp *amqp.Amqp
}

func NewMessageService(repo repository.Message, amqp *amqp.Amqp) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, message dto.MessageRequestDto) (int, error) {
	// TODO: send message to queue in transaction, if it failed, mark message as dead
	return s.repo.SaveMessage(ctx, message.ToModel())
}

func (s *MessageService) GetMessages(ctx context.Context) ([]model.Message, error) {
	return s.repo.GetMessages(ctx)
}