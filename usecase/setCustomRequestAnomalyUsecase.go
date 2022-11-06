package usecase

import (
	"context"
	"fmt"

	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type SetCustomRequestAnomalyUsecase struct {
	requestsAnomaliesService service.RequestsAnomaliesService
}

func NewSetCustomRequestAnomalyUsecase(requestsAnomaliesService service.RequestsAnomaliesService) *SetCustomRequestAnomalyUsecase {
	return &SetCustomRequestAnomalyUsecase{requestsAnomaliesService: requestsAnomaliesService}
}

func (uc *SetCustomRequestAnomalyUsecase) Execute(ctx context.Context, rootId string, isAnomaly bool) error {
	z := zap.S().With("context", "SetCustomRequestAnomalyUsecase")

	if err := uc.requestsAnomaliesService.SetCustomRequestAnomaly(ctx, rootId, isAnomaly); err != nil {
		z.Errorw("failed to upsert custom anomaly", "error", err.Error())
		return fmt.Errorf("failed to upsert custom anomaly: %w", err)
	}

	return nil
}
