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

func (r *MessageRepo) StartTx(ctx context.Context) (pgx.Tx, error) {
	return r.Pool.Begin(ctx)
}

func (r *MessageRepo) SaveMessage(ctx context.Context, message model.Message, tx pgx.Tx) (model.Message, error) {
	message.CreatedAt = time.Now()
	sql, args, _ := r.Builder.Insert("exchange.message").
		Columns("routing", "message", "dead", "retried", "created_at").
		Values(message.Routing, message.Message, message.Dead, message.Retried, message.CreatedAt).
		Suffix("RETURNING id").
		ToSql()

	if tx != nil {
		err := tx.QueryRow(ctx, sql, args...).Scan(&message.ID)
		if err != nil {
			return message, fmt.Errorf("MessageRepo.SaveMessage - tx.QueryRow: %v", err)
		}
	} else {
		err := r.Pool.QueryRow(ctx, sql, args...).Scan(&message.ID)
		if err != nil {
			return message, fmt.Errorf("MessageRepo.SaveMessage - pool.QueryRow: %v", err)
		}
	}

	return message, nil
}

func (r *MessageRepo) SaveMessages(ctx context.Context, messages []model.Message) error {
	batch := &pgx.Batch{}
	for _, msg := range messages {
		sql, args, _ := r.Builder.Insert("exchange.message").
			Columns("routing", "message", "dead", "retried", "created_at").
			Values(msg.Routing, msg.Message, msg.Dead, msg.Retried, time.Now()).
			ToSql()
		batch.Queue(sql, args...)
	}
	result := r.Pool.SendBatch(ctx, batch)
	_, err := result.Exec()
	if err != nil {
		return fmt.Errorf("MessageRepo.SaveMessages - r.Pool.SendBatch: %v", err)
	}
	err = result.Close()
	if err != nil {
		return fmt.Errorf("MessageRepo.SaveMessages - result.Close: %v", err)
	}
	return nil
}

type UpdateMessageInput struct {
	Routing *string            `json:"routing" db:"routing"`
	Message *model.AmqpMessage `json:"message" db:"message"`
	Dead    *bool              `json:"dead" db:"dead"`
	Retried *bool              `json:"retried" db:"retried"`
}

func (r *MessageRepo) UpdateMessage(ctx context.Context, ID int, message UpdateMessageInput, tx pgx.Tx) error {
	b := r.Builder.Update("exchange.message").Where("id = ?", ID)
	if message.Routing != nil {
		b.Set("routing", message.Routing)
	}
	if message.Message != nil {
		b.Set("message", message.Message)
	}
	if message.Dead != nil {
		b.Set("message", message.Message)
	}
	if message.Retried != nil {
		b.Set("message", message.Message)
	}

	sql, args, _ := b.ToSql()
	if tx != nil {
		fmt.Println("UPDATE SQL: ", sql)
		_, err := tx.Exec(ctx, sql, args)
		if err != nil {
			return err
		}

	} else {
		_, err := r.Pool.Exec(ctx, sql, args)
		if err != nil {
			return err
		}
	}
	return nil
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
			&message.UpdatedAt,
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
