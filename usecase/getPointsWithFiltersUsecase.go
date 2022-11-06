package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetPointsWithFiltersUsecase struct {
	hcsService service.HCSService
}

func NewGetPointsWithFiltersUsecase(hcsService service.HCSService) *GetPointsWithFiltersUsecase {
	return &GetPointsWithFiltersUsecase{hcsService: hcsService}
}

func (uc *GetPointsWithFiltersUsecase) Execute(
	ctx context.Context,
	filters *dto.GetPointsRequestQueryDTO,
) ([]entity.HCSPoint, error) {
	z := zap.S().With("context", "GetPointsWithFiltersUsecase")

	points, err := uc.hcsService.GetPointsWithFilters(ctx, filters)
	if err != nil {
		z.Errorw("failed to get points", "error", err.Error())
	}

	return points, err
}
