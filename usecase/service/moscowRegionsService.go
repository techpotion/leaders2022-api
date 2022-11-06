package service

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
)

type MoscowRegionsService interface {
	GetUniqueRegions(ctx context.Context) ([]string, error)
}

type moscowRegionsService struct {
	moscowRegionsRepository repository.MoscowRegionsRepository
}

func NewMoscowRegionsService(moscowRegionsRepository repository.MoscowRegionsRepository) *moscowRegionsService {
	return &moscowRegionsService{moscowRegionsRepository: moscowRegionsRepository}
}

func (s *moscowRegionsService) GetUniqueRegions(ctx context.Context) ([]string, error) {
	return s.moscowRegionsRepository.GetUniqueRegions(ctx)
}
