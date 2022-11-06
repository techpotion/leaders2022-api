package service

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
)

type HCSService interface {
	CountPointsWithFilters(ctx context.Context, filters *dto.CountPointsRequestDTO) (int, error)
	GetPointsWithFilters(ctx context.Context, filters *dto.GetPointsRequestQueryDTO) ([]entity.HCSPoint, error)
	GetRequestsByPointsIds(ctx context.Context, ids []string) ([]*entity.Request, error)
	GetPrimaryRequestDatetime(ctx context.Context, current *entity.Request) (*time.Time, error)
	CountRequestsByClosureTime(ctx context.Context, from, to time.Time) (int, error)
	GetRequestsByClosureTime(ctx context.Context, from, to time.Time, limit, offset int) ([]*entity.Request, error)
	CountRequestsFull(ctx context.Context, filters *dto.CountRequestsFullRequestDTO) (int, error)
	GetRequestsFull(ctx context.Context, filters *dto.GetRequestsFullRequestDTO) ([]*entity.RequestFull, error)
}

type hcsService struct {
	hcsRepository repository.HCSRepository
}

func NewHCSService(hcsRepository repository.HCSRepository) *hcsService {
	return &hcsService{hcsRepository: hcsRepository}
}

func (s *hcsService) GetPointsWithFilters(
	ctx context.Context,
	filters *dto.GetPointsRequestQueryDTO,
) ([]entity.HCSPoint, error) {
	return s.hcsRepository.GetPointsWithFilters(ctx, filters)
}

func (s *hcsService) CountPointsWithFilters(ctx context.Context, filters *dto.CountPointsRequestDTO) (int, error) {
	return s.hcsRepository.CountPointsWithFilters(ctx, filters)
}

func (s *hcsService) GetRequestsByPointsIds(ctx context.Context, ids []string) ([]*entity.Request, error) {
	return s.hcsRepository.GetRequestsByPointsIds(ctx, ids)
}

func (s *hcsService) GetPrimaryRequestDatetime(ctx context.Context, current *entity.Request) (*time.Time, error) {
	return s.hcsRepository.GetPrimaryRequestDatetime(ctx, current)
}

func (s *hcsService) CountRequestsByClosureTime(ctx context.Context, from, to time.Time) (int, error) {
	return s.hcsRepository.CountRequestsByClosureTime(ctx, from, to)
}

func (s *hcsService) GetRequestsByClosureTime(ctx context.Context, from, to time.Time, limit, offset int) ([]*entity.Request, error) {
	return s.hcsRepository.GetRequestsByClosureTime(ctx, from, to, limit, offset)
}

func (s *hcsService) CountRequestsFull(ctx context.Context, filters *dto.CountRequestsFullRequestDTO) (int, error) {
	if len(filters.AnomalyCases) == 0 {
		filters.AnomalyCases = make([]int, 0, 1)
	}

	return s.hcsRepository.CountRequestsFull(ctx, filters)
}

func (s *hcsService) GetRequestsFull(ctx context.Context, filters *dto.GetRequestsFullRequestDTO) ([]*entity.RequestFull, error) {
	if len(filters.AnomalyCases) == 0 {
		filters.AnomalyCases = make([]int, 0, 1)
	}

	return s.hcsRepository.GetRequestsFull(ctx, filters)
}
