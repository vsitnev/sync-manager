package pgdb

import (
	"context"
	"fmt"

	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/pkg/postgres"
)

type MessageRepo struct {
	*postgres.Postgres
}

func NewMessageRepo(db *postgres.Postgres) *MessageRepo {
	return &MessageRepo{db}
}

func (r *MessageRepo) SaveMessage(ctx context.Context, message model.Message) (int, error) {
	sql, args, _ := r.Builder.Insert("exchange.message").
	Columns("routing", "message", "dead").
	Values(message.Routing, message.Message, message.Dead).
	Suffix("RETURNING id").
	ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("MessageRepo.SaveMessage - row.Scan: %v", err)
	}
	return id, nil
}

func (r *MessageRepo) GetMessages(ctx context.Context) ([]model.Message, error) {
	var messages []model.Message

	sql, args, _ := r.Builder.Select("*").From("exchange.message").ToSql()
	
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("LocationRepo.GetLocations - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message model.Message
		err = rows.Scan(
			&message.ID,
			&message.Routing,
			&message.Message,
			&message.Dead,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("LocationRepo.GetLocations - rows.Scan: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}