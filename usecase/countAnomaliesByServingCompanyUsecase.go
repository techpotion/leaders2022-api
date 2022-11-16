package usecase

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountAnomaliesByServingCompaniesUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewCountAnomaliesByServingCompaniesUsecase(
	requestsAnomaliesService service.RequestsAnomaliesService,
) *CountAnomaliesByServingCompaniesUsecase {
	return &CountAnomaliesByServingCompaniesUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *CountAnomaliesByServingCompaniesUsecase) Execute(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByServingCompany, error) {
	z := zap.S().With("context", "CountAnomaliesByServingCompaniesUsecase")

	c, err := uc.requestsAnomaliesService.CountAnomaliesByServingCompanies(ctx, dateFrom, dateTo, region)
	if err != nil {
		z.Errorw(
			"failed to count groupped anomalies by serving companies",
			"error", err.Error(),
			"date_from", dateFrom,
			"dateTo", dateTo,
			"region", region,
		)
	}

	return c, err
}
