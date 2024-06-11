package service

import (
	"context"

	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/pkg/amqpclient"
)

type MessageService struct {
	repo repository.Message
	amqp *amqpclient.Amqp
}

func NewMessageService(repo repository.Message, amqp *amqpclient.Amqp) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, message model.Message) (int, error) {
	// TODO: send message to queue in transaction, if it failed, mark message as dead
	return s.repo.SaveMessage(ctx, message)
}

func (s *MessageService) GetMessages(ctx context.Context, input MessageInput) ([]model.Message, error) {
	return s.repo.GetMessagesPagination(ctx, input.Source, input.Routing, input.SortType, input.Limit, input.Offset)
}

func (s *MessageService) GetMessageByID(ctx context.Context, ID int) (model.Message, error) {
	return s.repo.GetMessageByID(ctx, ID)
}
