package usecase

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountAnomaliesByOwnerCompaniesUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewCountAnomaliesByOwnerCompaniesUsecase(
	requestsAnomaliesService service.RequestsAnomaliesService,
) *CountAnomaliesByOwnerCompaniesUsecase {
	return &CountAnomaliesByOwnerCompaniesUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *CountAnomaliesByOwnerCompaniesUsecase) Execute(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByOwnerCompany, error) {
	z := zap.S().With("context", "CountAnomaliesByOwnerCompaniesUsecase")

	c, err := uc.requestsAnomaliesService.CountAnomaliesByOwnerCompanies(ctx, dateFrom, dateTo, region)
	if err != nil {
		z.Errorw(
			"failed to count anomalies groupped by owner companies",
			"error", err.Error(),
			"date_from", dateFrom,
			"dateTo", dateTo,
			"region", region,
		)
	}

	return c, err
}
