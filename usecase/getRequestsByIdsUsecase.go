package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetRequestsByIdsUsecase struct {
	hcsService service.HCSService
}

func NewGetRequestsByIdsUsecase(hcsService service.HCSService) *GetRequestsByIdsUsecase {
	return &GetRequestsByIdsUsecase{hcsService: hcsService}
}

func (uc *GetRequestsByIdsUsecase) Execute(ctx context.Context, ids []string) ([]*entity.Request, error) {
	z := zap.S().With("context", "GetRequestsByIdsUsecase")

	points, err := uc.hcsService.GetRequestsByPointsIds(ctx, ids)
	if err != nil {
		z.Errorw("failed to get requests by points", "error", err.Error())
	}

	return points, err
}
