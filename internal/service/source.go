package service

import (
	"context"
	"github.com/vsitnev/sync-manager/internal/dto"
	"github.com/vsitnev/sync-manager/internal/repository"
	"github.com/vsitnev/sync-manager/internal/repository/pgdb"
)

type SourceService struct {
	repo repository.Source
}

func NewSourceService(repo repository.Source) *SourceService {
	return &SourceService{
		repo: repo,
	}
}

func (s *SourceService) SaveSource(ctx context.Context, source dto.Source) (int, error) {
	return s.repo.SaveSource(ctx, source)
}

func (s *SourceService) GetSources(ctx context.Context, filter PaginationFilter) ([]dto.Source, error) {
	var sources []dto.Source
	data, err := s.repo.GetSourcesPagination(ctx, filter.SortType, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	for i := range data {
		sources = append(sources, data[i].ToDto())
	}

	return sources, nil
}

func (s *SourceService) GetSourceByID(ctx context.Context, ID int) (dto.Source, error) {
	var source dto.Source
	data, err := s.repo.GetSourceByID(ctx, ID)
	if err != nil {
		return source, err
	}
	return data.ToDto(), nil
}

func (s *SourceService) GetSourceByCode(ctx context.Context, code string) (dto.Source, error) {
	var source dto.Source
	data, err := s.repo.GetSourceByCode(ctx, code)
	if err != nil {
		return source, err
	}
	return data.ToDto(), nil
}

func (s *SourceService) UpdateSource(ctx context.Context, ID int, source pgdb.UpdateSourceInput) error {
	return s.repo.UpdateSource(ctx, ID, source)
}
