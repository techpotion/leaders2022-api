package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetRequestsFullUsecase struct {
	hcsService service.HCSService
}

func NewGetRequestsFullUsecase(hcsService service.HCSService) *GetRequestsFullUsecase {
	return &GetRequestsFullUsecase{hcsService: hcsService}
}

func (uc *GetRequestsFullUsecase) Execute(
	ctx context.Context,
	filters *dto.GetRequestsFullRequestDTO,
) ([]*entity.RequestFull, error) {
	z := zap.S().With("context", "GetRequestsFullUsecase")

	requests, err := uc.hcsService.GetRequestsFull(ctx, filters)
	if err != nil {
		z.Errorw("failed to get points", "error", err.Error())
	}

	return requests, err
}
