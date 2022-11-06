package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetUniqueMoscowRegionsUsecase struct {
	regionsService service.MoscowRegionsService
}

func NewGetUniqueMoscowRegionsUsecase(regionsService service.MoscowRegionsService) *GetUniqueMoscowRegionsUsecase {
	return &GetUniqueMoscowRegionsUsecase{regionsService: regionsService}
}

func (uc *GetUniqueMoscowRegionsUsecase) Execute(ctx context.Context) ([]string, error) {
	z := zap.S().With("context", "GetUniqueMoscowRegionsUsecase")

	regions, err := uc.regionsService.GetUniqueRegions(ctx)
	if err != nil {
		z.Errorw("failed to get unique regions", "error", err.Error())
	}

	return regions, err
}
