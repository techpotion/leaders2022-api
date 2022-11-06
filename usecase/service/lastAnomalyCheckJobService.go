package service

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
)

type LastAnomalyCheckJobService interface {
	GetLastTimestamp(ctx context.Context) (time.Time, bool, error)
	SetJobAsActive(ctx context.Context) (time.Time, error)
	SetJobAsInactive(ctx context.Context, newLastJobTime *time.Time) (time.Time, error)
}

type lastAnomalyCheckJobService struct {
	lastAnomalyCheckJobRepository repository.LastAnomalyCheckJobRepository
}

func NewLastAnomalyCheckJobService(
	lastAnomalyCheckJobRepository repository.LastAnomalyCheckJobRepository,
) *lastAnomalyCheckJobService {
	return &lastAnomalyCheckJobService{lastAnomalyCheckJobRepository: lastAnomalyCheckJobRepository}
}

func (s *lastAnomalyCheckJobService) GetLastTimestamp(ctx context.Context) (time.Time, bool, error) {
	return s.lastAnomalyCheckJobRepository.GetLastTimestamp(ctx)
}

func (s *lastAnomalyCheckJobService) SetJobAsActive(ctx context.Context) (time.Time, error) {
	return s.lastAnomalyCheckJobRepository.SetJobAsActive(ctx)
}

func (s *lastAnomalyCheckJobService) SetJobAsInactive(ctx context.Context, newLastJobTime *time.Time) (time.Time, error) {
	return s.lastAnomalyCheckJobRepository.SetJobAsInactive(ctx, newLastJobTime)
}
