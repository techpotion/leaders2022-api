package usecase

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountAnomaliesGrouppedUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewCountAnomaliesGrouppedUsecase(requestsAnomaliesService service.RequestsAnomaliesService) *CountAnomaliesGrouppedUsecase {
	return &CountAnomaliesGrouppedUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *CountAnomaliesGrouppedUsecase) Execute(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.CountAnomaliesGroupped, error) {
	z := zap.S().With("context", "CountAnomaliesGrouppedUsecase")

	m, err := uc.requestsAnomaliesService.CountAnomaliesGroupped(ctx, dateFrom, dateTo, region)
	if err != nil {
		z.Errorw(
			"failed to count groupped anomalies with percent",
			"error", err.Error(),
			"date_from", dateFrom,
			"dateTo", dateTo,
			"region", region,
		)
	}

	return m, err
}
