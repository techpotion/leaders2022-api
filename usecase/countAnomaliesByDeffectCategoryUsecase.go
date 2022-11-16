package usecase

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type CountAnomaliesByDeffectCategoriesUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewCountAnomaliesByDeffectCategoriesUsecase(
	requestsAnomaliesService service.RequestsAnomaliesService,
) *CountAnomaliesByDeffectCategoriesUsecase {
	return &CountAnomaliesByDeffectCategoriesUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *CountAnomaliesByDeffectCategoriesUsecase) Execute(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByDeffectCategory, error) {
	z := zap.S().With("context", "CountAnomaliesByDeffectCategoriesUsecase")

	c, err := uc.requestsAnomaliesService.CountAnomaliesByDeffectCategories(ctx, dateFrom, dateTo, region)
	if err != nil {
		z.Errorw(
			"failed to count anomalies groupped by deffect categories",
			"error", err.Error(),
			"date_from", dateFrom,
			"dateTo", dateTo,
			"region", region,
		)
	}

	return c, err
}
