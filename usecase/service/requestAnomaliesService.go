package service

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
)

type RequestsAnomaliesService interface {
	UpsertAnomaly(ctx context.Context, a *entity.RequestAnomaly) error
	UpsertAnomalies(ctx context.Context, as []*entity.RequestAnomaly) error
	GetAnomaliesByRootIds(ctx context.Context, ids []string, cases []int) ([]entity.RequestAnomaly, error)
	SetCustomRequestAnomaly(ctx context.Context, rootId string, isAnomaly bool) error
	CountAnomalies(ctx context.Context, dateFrom, dateTo time.Time, region string) (int, float32, error)
	CountAnomaliesGroupped(ctx context.Context, dateFrom, dateTo time.Time, region string) ([]*entity.CountAnomaliesGroupped, error)
	GetAnomaliesAmountDynamics(ctx context.Context, dateFrom, dateTo time.Time, region string) ([]*entity.AnomaliesDynamics, error)
	CountAnomaliesByOwnerCompanies(
		ctx context.Context,
		dateFrom, dateTo time.Time,
		region string,
	) ([]*entity.AnomaliesByOwnerCompany, error)
	CountAnomaliesByServingCompanies(
		ctx context.Context,
		dateFrom, dateTo time.Time,
		region string,
	) ([]*entity.AnomaliesByServingCompany, error)
	CountAnomaliesByDeffectCategories(
		ctx context.Context,
		dateFrom, dateTo time.Time,
		region string,
	) ([]*entity.AnomaliesByDeffectCategory, error)
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

func (s *requestsAnomaliesService) CountAnomalies(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) (count int, perc float32, err error) {
	return s.requestsAnomaliesRepository.CountAnomalies(ctx, dateFrom, dateTo, region)
}

func (s *requestsAnomaliesService) CountAnomaliesGroupped(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.CountAnomaliesGroupped, error) {
	return s.requestsAnomaliesRepository.CountAnomaliesGroupped(ctx, dateFrom, dateTo, region)
}

func (s *requestsAnomaliesService) GetAnomaliesAmountDynamics(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesDynamics, error) {
	return s.requestsAnomaliesRepository.GetAnomaliesAmountDynamics(ctx, dateFrom, dateTo, region)
}

func (s *requestsAnomaliesService) CountAnomaliesByOwnerCompanies(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByOwnerCompany, error) {
	return s.requestsAnomaliesRepository.CountAnomaliesByOwnerCompanies(ctx, dateFrom, dateTo, region)
}

func (s *requestsAnomaliesService) CountAnomaliesByServingCompanies(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByServingCompany, error) {
	return s.requestsAnomaliesRepository.CountAnomaliesByServingCompanies(ctx, dateFrom, dateTo, region)
}

func (s *requestsAnomaliesService) CountAnomaliesByDeffectCategories(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByDeffectCategory, error) {
	return s.requestsAnomaliesRepository.CountAnomaliesByDeffectCategories(ctx, dateFrom, dateTo, region)
}
