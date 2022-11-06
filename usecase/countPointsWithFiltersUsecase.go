package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountPointsWithFiltersUsecase struct {
	hcsService service.HCSService
}

func NewCountPointsWithFiltersUsecase(hcsService service.HCSService) *CountPointsWithFiltersUsecase {
	return &CountPointsWithFiltersUsecase{hcsService: hcsService}
}

func (uc *CountPointsWithFiltersUsecase) Execute(
	ctx context.Context,
	filters *dto.CountPointsRequestDTO,
) (int, error) {
	z := zap.S().With("context", "CountPointsWithFiltersUsecase")

	count, err := uc.hcsService.CountPointsWithFilters(ctx, filters)
	if err != nil {
		z.Errorw("failed to get points", "error", err.Error())
	}

	return count, err
}
