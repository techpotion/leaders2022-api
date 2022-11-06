package usecase

import (
	"context"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetRequestsAnomaliesByIdsUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewGetRequestsAnomaliesByIdsUsecase(requestsAnomaliesService service.RequestsAnomaliesService) *GetRequestsAnomaliesByIdsUsecase {
	return &GetRequestsAnomaliesByIdsUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *GetRequestsAnomaliesByIdsUsecase) Execute(ctx context.Context, ids []string, cases []int) ([]entity.RequestAnomaly, error) {
	z := zap.S().With("context", "GetRequestsAnomaliesByIdsUsecase")

	anomalies, err := uc.requestsAnomaliesService.GetAnomaliesByRootIds(ctx, ids, cases)
	if err != nil {
		z.Errorw("failed to get requests by points", "error", err.Error())
	}

	return anomalies, err
}
