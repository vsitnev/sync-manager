package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository/repoerr"
	"github.com/vsitnev/sync-manager/pkg/postgres"
	"time"
)

const (
	maxPaginationLimit     = 10
	defaultPaginationLimit = 10

	DateSortType string = "date"
)

type MessageRepo struct {
	*postgres.Postgres
}

func NewMessageRepo(db *postgres.Postgres) *MessageRepo {
	return &MessageRepo{db}
}

func (r *MessageRepo) SaveMessage(ctx context.Context, message model.Message) (int, error) {
	sql, args, _ := r.Builder.Insert("exchange.message").
		Columns("routing", "message", "dead", "retried", "created_at").
		Values(message.Routing, message.Message, message.Dead, message.Retried, time.Now()).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("MessageRepo.SaveMessage - row.Scan: %v", err)
	}
	return id, nil
}

func (r *MessageRepo) GetMessagesPagination(ctx context.Context, source, routing string, sortType string, limit int, offset int) ([]model.Message, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}
	if limit == 0 {
		limit = defaultPaginationLimit
	}

	var orderBySql string
	switch sortType {
	case "":
		orderBySql = "created_at DESC"
	case DateSortType:
		orderBySql = "created_at DESC"
	default:
		return nil, fmt.Errorf("MessageRepo.GetMessagesPagination: unknown sort type - %s", sortType)
	}

	var whereClauses []squirrel.Sqlizer
	if source != "" {
		whereClauses = append(whereClauses, squirrel.Eq{"message ->> 'source'": source})
	}
	if routing != "" {
		whereClauses = append(whereClauses, squirrel.Eq{"routing": routing})
	}

	query := r.Builder.Select("*").From("exchange.message").
		OrderBy(orderBySql).
		Limit(uint64(limit)).
		Offset(uint64(offset))
	if len(whereClauses) > 0 {
		query = query.Where(squirrel.And(whereClauses))
	}
	sql, args, _ := query.ToSql()

	fmt.Println(sql)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("LocationRepo.GetLocations - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		err = rows.Scan(
			&message.ID,
			&message.Routing,
			&message.Message,
			&message.Dead,
			&message.Retried,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("LocationRepo.GetLocations - rows.Scan: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepo) GetMessageByID(ctx context.Context, ID int) (model.Message, error) {
	var message model.Message
	sql, args, _ := r.Builder.Select("*").From("exchange.message").
		Where("id = ?", ID).ToSql()

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&message.ID,
		&message.Routing,
		&message.Message,
		&message.Dead,
		&message.Retried,
		&message.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return message, repoerr.ErrNotFound
		}
		return message, fmt.Errorf("MessageRepo.GetMessageByID - tx.QueryRow: %v", err)
	}

	return message, nil
}
