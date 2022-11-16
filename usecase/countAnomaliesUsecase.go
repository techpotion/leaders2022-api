package usecase

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountAnomaliesUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewCountAnomaliesUsecase(requestsAnomaliesService service.RequestsAnomaliesService) *CountAnomaliesUsecase {
	return &CountAnomaliesUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *CountAnomaliesUsecase) Execute(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) (count int, perc float32, err error) {
	z := zap.S().With("context", "CountAnomaliesUsecase")

	count, perc, err = uc.requestsAnomaliesService.CountAnomalies(ctx, dateFrom, dateTo, region)
	if err != nil {
		z.Errorw(
			"failed to count anomalies with percent",
			"error", err.Error(),
			"date_from", dateFrom,
			"dateTo", dateTo,
			"region", region,
		)
	}

	return count, perc, err
}
