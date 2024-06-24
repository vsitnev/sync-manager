package service

import (
	"context"
	"github.com/vsitnev/sync-manager/internal/dto"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/internal/repository/pgdb"
	"github.com/vsitnev/sync-manager/pkg/amqpclient"
)

type MessageInput struct {
	Source  string
	Routing string
	PaginationFilter
}

type PaginationFilter struct {
	SortType string
	Offset   int
	Limit    int
}

type Message interface {
	SendMessage(ctx context.Context, message model.Message) (SendMessageResponse, error)

	SaveMessages(ctx context.Context, messages []model.Message) error
	GetMessages(ctx context.Context, input MessageInput) ([]dto.Message, error)
	GetMessageByID(ctx context.Context, ID int) (dto.Message, error)
}

type Source interface {
	SaveSource(ctx context.Context, source dto.Source) (int, error)
	GetSources(ctx context.Context, filter PaginationFilter) ([]dto.Source, error)
	GetSourceByID(ctx context.Context, ID int) (dto.Source, error)
	GetSourceByCode(ctx context.Context, code string) (dto.Source, error)
	UpdateSource(ctx context.Context, ID int, source pgdb.UpdateSourceInput) error
}

type ServiceDeps struct {
	Reps *repository.Repositories
	Amqp *amqpclient.Amqp
}

type Services struct {
	Message Message
	Source  Source
}

func NewServices(deps ServiceDeps) *Services {
	sourceService := NewSourceService(deps.Reps)
	return &Services{
		Message: NewMessageService(deps.Reps, sourceService, deps.Amqp),
		Source:  sourceService,
	}
}
