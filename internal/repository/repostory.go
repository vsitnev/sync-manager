package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/vsitnev/sync-manager/internal/dto"

	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository/pgdb"
	"github.com/vsitnev/sync-manager/pkg/postgres"
)

type Message interface {
	StartTx(ctx context.Context) (pgx.Tx, error)

	SaveMessage(ctx context.Context, message model.Message, tx pgx.Tx) (model.Message, error)
	SaveMessages(ctx context.Context, messages []model.Message) error
	UpdateMessage(ctx context.Context, ID int, message pgdb.UpdateMessageInput, tx pgx.Tx) error
	GetMessagesPagination(ctx context.Context, source, routing string, sortType string, offset int, limit int) ([]model.Message, error)
	GetMessageByID(ctx context.Context, ID int) (model.Message, error)
}

type Source interface {
	SaveSource(ctx context.Context, source dto.Source) (int, error)
	GetSourcesPagination(ctx context.Context, sortType string, limit int, offset int) ([]model.Source, error)
	GetSourceByID(ctx context.Context, ID int) (model.Source, error)
	GetSourceByCode(ctx context.Context, ID string) (model.Source, error)
	UpdateSource(ctx context.Context, ID int, source pgdb.UpdateSourceInput) error
}

type Repositories struct {
	Message
	Source
}

func NewRepositories(db *postgres.Postgres) *Repositories {
	return &Repositories{
		Message: pgdb.NewMessageRepo(db),
		Source:  pgdb.NewSourceRepo(db),
	}
}
