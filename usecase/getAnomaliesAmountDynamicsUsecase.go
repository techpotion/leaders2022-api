package usecase

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetAnomaliesAmountDynamicsUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewGetAnomaliesAmountDynamicsUsecase(requestsAnomaliesService service.RequestsAnomaliesService) *GetAnomaliesAmountDynamicsUsecase {
	return &GetAnomaliesAmountDynamicsUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *GetAnomaliesAmountDynamicsUsecase) Execute(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesDynamics, error) {
	z := zap.S().With("context", "GetAnomaliesAmountDynamicsUsecase")

	aads, err := uc.requestsAnomaliesService.GetAnomaliesAmountDynamics(ctx, dateFrom, dateTo, region)
	if err != nil {
		z.Errorw(
			"failed to get anomalies amount dynamics",
			"error", err.Error(),
			"date_from", dateFrom,
			"dateTo", dateTo,
			"region", region,
		)
	}

	return aads, err
}
