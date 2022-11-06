package service

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
)

type RequestsAnomaliesService interface {
	UpsertAnomaly(ctx context.Context, a *entity.RequestAnomaly) error
	UpsertAnomalies(ctx context.Context, as []*entity.RequestAnomaly) error
	GetAnomaliesByRootIds(ctx context.Context, ids []string, cases []int) ([]entity.RequestAnomaly, error)
	SetCustomRequestAnomaly(ctx context.Context, rootId string, isAnomaly bool) error
}

type requestsAnomaliesService struct {
	requestsAnomaliesRepository repository.RequestsAnomaliesRepository
}

func NewRequestsAnomaliesService(requestsAnomaliesRepository repository.RequestsAnomaliesRepository) *requestsAnomaliesService {
	return &requestsAnomaliesService{requestsAnomaliesRepository: requestsAnomaliesRepository}
}

func (s *requestsAnomaliesService) UpsertAnomaly(ctx context.Context, a *entity.RequestAnomaly) error {
	return s.requestsAnomaliesRepository.UpsertAnomaly(ctx, a)
}

func (s *requestsAnomaliesService) UpsertAnomalies(ctx context.Context, as []*entity.RequestAnomaly) error {
	return s.requestsAnomaliesRepository.UpsertAnomalies(ctx, as)
}

func (s *requestsAnomaliesService) GetAnomaliesByRootIds(ctx context.Context, ids []string, cases []int) ([]entity.RequestAnomaly, error) {
	if len(cases) == 0 {
		cases = make([]int, 0, 1)
	}

	return s.requestsAnomaliesRepository.GetAnomaliesByRootIds(ctx, ids, cases)
}

func (s *requestsAnomaliesService) SetCustomRequestAnomaly(ctx context.Context, rootId string, isAnomaly bool) error {
	return s.requestsAnomaliesRepository.SetCustomRequestAnomaly(ctx, rootId, isAnomaly)
}
