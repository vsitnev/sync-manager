package service

import (
	"context"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
)

type AmqpMessageService struct {
	repo repository.Message
}

func NewAmqpMessageService(repo repository.Message) *AmqpMessageService {
	return &AmqpMessageService{
		repo: repo,
	}
}

func (s *AmqpMessageService) SaveFromAmqp(ctx context.Context, message model.AmqpMessage) error {
	return nil
}
