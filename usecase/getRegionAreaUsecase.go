package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetRegionAreaUsecase struct {
	hcsService service.HCSService
}

func NewGetRegionAreaUsecase(hcsService service.HCSService) *GetRegionAreaUsecase {
	return &GetRegionAreaUsecase{hcsService: hcsService}
}

func (uc *GetRegionAreaUsecase) Execute(ctx context.Context, region string) (string, error) {
	z := zap.S().With("context", "GetRegionAreaUsecase")

	area, err := uc.hcsService.GetRegionArea(ctx, region)
	if err != nil {
		z.Errorw("failed to get region area", "error", err.Error(), "region", region)
		return "", err
	}

	return area, nil
}
