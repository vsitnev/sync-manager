package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/vsitnev/sync-manager/internal/dto"
	"github.com/vsitnev/sync-manager/internal/model"
	"github.com/vsitnev/sync-manager/internal/repository/repoerr"
	"github.com/vsitnev/sync-manager/pkg/postgres"
	"time"
)

type SourceRepo struct {
	*postgres.Postgres
}

func NewSourceRepo(db *postgres.Postgres) *SourceRepo {
	return &SourceRepo{db}
}

func (r *SourceRepo) SaveSource(ctx context.Context, source dto.Source) (int, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("SourceRepo.SaveSource - r.Pool.BeginTx: %v", err)
	}
	defer func() { tx.Rollback(ctx) }()

	var id int
	sql, args, _ := r.Builder.Insert("exchange.source").
		Columns("name", "description", "code", "receive_method", "created_at").
		Values(source.Name, source.Description, source.Code, source.ReceiveMethod, time.Now()).
		Suffix("RETURNING id").
		ToSql()
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("SourceRepo.SaveSource - r.Pool.QueryRow: %v", err)
	}

	if source.Routes != nil && len(source.Routes) != 0 {
		err = r.saveRoutes(ctx, id, source.Routes, tx)
		if err != nil {
			return 0, fmt.Errorf("SourceRepo.SaveSource - r.SaveRoutes: %v", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("SourceRepo.SaveSource - tx.Commit: %v", err)
	}
	return id, nil
}

func (r *SourceRepo) saveRoutes(ctx context.Context, sourceID int, routes []dto.Route, tx pgx.Tx) error {
	batch := &pgx.Batch{}
	for _, route := range routes {
		sql, args, _ := r.Builder.Insert("exchange.route").
			Columns("name", "url", "source_fk").
			Values(route.Name, route.Url, sourceID).
			ToSql()
		batch.Queue(sql, args...)
	}
	result := tx.SendBatch(ctx, batch)
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

func (r *SourceRepo) GetSourcesPagination(ctx context.Context, sortType string, limit int, offset int) ([]model.Source, error) {
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
		return nil, fmt.Errorf("SourceRepo.GetSourcesPagination: unknown sort type - %s", sortType)
	}

	sql, args, err := r.Builder.Select("*").From("exchange.source").
		OrderBy(orderBySql).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SourceRepo.GetSourcesPagination - r.Builder.ToSql: %v", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SourceRepo.GetSourcesPagination - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var sources []model.Source
	for rows.Next() {
		var source model.Source
		err = rows.Scan(
			&source.ID,
			&source.Name,
			&source.Description,
			&source.Code,
			&source.ReceiveMethod,
			&source.CreatedAt,
			&source.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("SourceRepo.GetSourcesPagination - rows.Scan: %v", err)
		}
		sources = append(sources, source)
	}

	for i, _ := range sources {
		sources[i].Routes, err = r.getRouteBySourceID(ctx, sources[i].ID)
		if err != nil {
			return nil, fmt.Errorf("SourceRepo.GetSourcesPagination - r.getRouteBySourceID: %v", err)
		}
	}

	return sources, nil
}

func (r *SourceRepo) getRouteBySourceID(ctx context.Context, sourceID int) ([]model.Route, error) {
	sql, args, err := r.Builder.Select("*").From("exchange.route").
		Where("source_fk = ?", sourceID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Builder.ToSql: %v", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var routes []model.Route
	for rows.Next() {
		var route model.Route
		err = rows.Scan(
			&route.ID,
			&route.Name,
			&route.Url,
			&route.SourceID,
		)
		if err != nil {
			return nil, fmt.Errorf("SourceRepo.GetRouteBySourceID - rows.Scan: %v", err)
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func (r *SourceRepo) GetSourceByID(ctx context.Context, ID int) (model.Source, error) {
	var source model.Source
	sql, args, _ := r.Builder.Select("*").From("exchange.source").
		Where("id = ?", ID).
		ToSql()
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&source.ID,
		&source.Name,
		&source.Description,
		&source.Code,
		&source.ReceiveMethod,
		&source.CreatedAt,
		&source.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return source, repoerr.ErrNotFound
		}
		return source, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.QueryRow: %v", err)
	}

	source.Routes, err = r.getRouteBySourceID(ctx, ID)
	if err != nil {
		return source, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.QueryRow: %v", err)
	}

	return source, nil
}

func (r *SourceRepo) GetSourceByCode(ctx context.Context, code string) (model.Source, error) {
	var source model.Source
	sql, args, _ := r.Builder.Select("*").From("exchange.source").
		Where("code = ?", code).
		ToSql()
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&source.ID,
		&source.Name,
		&source.Description,
		&source.Code,
		&source.ReceiveMethod,
		&source.CreatedAt,
		&source.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return source, repoerr.ErrNotFound
		}
		return source, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.QueryRow: %v", err)
	}

	source.Routes, err = r.getRouteBySourceID(ctx, source.ID)
	if err != nil {
		return source, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.QueryRow: %v", err)
	}

	return source, nil
}

type UpdateSourceInput struct {
	Name          *string
	Description   *string
	Code          *string
	ReceiveMethod *string
	Routes        *[]dto.Route
}

func (r *SourceRepo) UpdateSource(ctx context.Context, ID int, source UpdateSourceInput) error {
	if _, err := r.GetSourceByID(ctx, ID); err != nil {
		return err
	}

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("SourceRepo.UpdateSource - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	b := r.Builder.Update("exchange.source").Where("id = ?", ID)
	if source.Name != nil {
		b = b.Set("name", source.Name)
	}
	if source.Description != nil {
		b = b.Set("description", source.Description)
	}
	if source.Code != nil {
		b = b.Set("code", source.Code)
	}
	if source.ReceiveMethod != nil {
		b = b.Set("receive_method", source.ReceiveMethod)
	}
	sql, argsq, _ := b.ToSql()

	_, err = tx.Exec(ctx, sql, argsq...)
	if err != nil {
		return fmt.Errorf("SourceRepo.UpdateSource - tx.Exec: %v", err)
	}

	if source.Routes != nil {
		err = r.updateRoute(ctx, ID, *source.Routes, tx)
		if err != nil {
			return fmt.Errorf("SourceRepo.UpdateSource - r.updateRoute: %v", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("SourceRepo.UpdateSource - tx.Commit: %v", err)
	}
	fmt.Println("WE ARE HERE")
	return nil
}

type SaveRouteInput struct {
	Name *string
	Url  *string
}

func (r *SourceRepo) updateRoute(ctx context.Context, sourceID int, routes []dto.Route, tx pgx.Tx) error {
	sql, args, _ := r.Builder.Delete("exchange.route").Where("source_fk = ?", sourceID).ToSql()
	_, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SourceRepo.updateRoute - r.Pool.Exc: %w", err)
	}

	if routes != nil && len(routes) != 0 {
		err = r.saveRoutes(ctx, sourceID, routes, tx)
		if err != nil {
			return fmt.Errorf("SourceRepo.updateRoute - r.saveRoutes: %w", err)
		}
	}
	return nil
}
