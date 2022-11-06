package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountRequestsFullUsecase struct {
	hcsService service.HCSService
}

func NewCountRequestsFullUsecase(hcsService service.HCSService) *CountRequestsFullUsecase {
	return &CountRequestsFullUsecase{hcsService: hcsService}
}

func (uc *CountRequestsFullUsecase) Execute(
	ctx context.Context,
	filters *dto.CountRequestsFullRequestDTO,
) (int, error) {
	z := zap.S().With("context", "CountRequestsFullUsecase")

	count, err := uc.hcsService.CountRequestsFull(ctx, filters)
	if err != nil {
		z.Errorw("failed to get points", "error", err.Error())
	}

	return count, err
}
