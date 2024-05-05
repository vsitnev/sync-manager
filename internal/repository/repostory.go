package repository

import (
	"context"

	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository/pgdb"
	"github.com/vsitnev/sync-manager/pkg/postgres"
)


type Message interface {
	SaveMessage(ctx context.Context, message model.Message) (int, error)
	GetMessages(ctx context.Context) ([]model.Message, error)
}


type Repositories struct {
	Message
}

func NewRepositories(db *postgres.Postgres) *Repositories {
	return &Repositories{
		Message: pgdb.NewMessageRepo(db),
	}
}