package usecase

import (
	"context"
	"math"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type EnrichRequestsWithAnomaliesUsecase struct {
	lastAnomalyCheckJobService service.LastAnomalyCheckJobService
	hcsService                 service.HCSService
	modelService               service.ModelService
	requestsAnomaliesService   service.RequestsAnomaliesService
}

func NewEnrichRequestsWithAnomaliesUsecase(
	lastAnomalyCheckJobService service.LastAnomalyCheckJobService,
	hcsService service.HCSService,
	modelService service.ModelService,
	requestsAnomaliesService service.RequestsAnomaliesService,
) *EnrichRequestsWithAnomaliesUsecase {
	return &EnrichRequestsWithAnomaliesUsecase{
		hcsService:                 hcsService,
		lastAnomalyCheckJobService: lastAnomalyCheckJobService,
		modelService:               modelService,
		requestsAnomaliesService:   requestsAnomaliesService,
	}
}

const maxPaginationLimit = 20000

func (uc *EnrichRequestsWithAnomaliesUsecase) Execute(ctx context.Context) (alreadyActive bool) {
	now := time.Now()
	t2 := &now

	z := zap.S().With("context", "EnrichRequestsWithAnomaliesUsecase", "start_time", t2)

	_, isActive, err := uc.lastAnomalyCheckJobService.GetLastTimestamp(ctx)
	if err != nil {
		z.Errorw(
			"Failed to get anomaly check job info",
			"error", err.Error(),
		)

		return
	}

	if isActive {
		return true
	}

	lastJobTime, err := uc.lastAnomalyCheckJobService.SetJobAsActive(ctx)
	if err != nil {
		z.Errorw(
			"Failed to set job as active",
			"error", err.Error(),
		)

		return
	}

	defer func() {
		lastJobTime, err = uc.lastAnomalyCheckJobService.SetJobAsInactive(ctx, t2)
		if err != nil {
			z.Errorw(
				"Failed to set job as inactive",
				"error", err.Error(),
			)
		}
	}()

	reqCount, err := uc.hcsService.CountRequestsByClosureTime(ctx, lastJobTime, *t2)
	if err != nil {
		z.Errorw(
			"Failed to count requests by closure time",
			"error", err.Error(),
		)

		t2 = nil

		return
	}

	if reqCount == 0 {
		z.Infow("New requests were not found")
		return
	}

	z.Infow("Found new requests", "count", reqCount)

	for i := 0; i < int(math.Ceil(float64(reqCount)/float64(maxPaginationLimit))); i++ {
		offset := maxPaginationLimit * i

		reqs, err := uc.hcsService.GetRequestsByClosureTime(ctx, lastJobTime, *t2, maxPaginationLimit, offset)
		if err != nil {
			z.Errorw(
				"Failed to get requests",
				"error", err.Error(),
			)

			t2 = nil

			break
		}

		modelAnomalyPredictions, err := uc.modelService.GetPredictions(ctx, uc.requestEntitiesToDtos(reqs))
		if err != nil {
			z.Errorw(
				"Failed to get model anomalies predictions",
				"error", err.Error(),
			)

			t2 = nil

			break
		}

		if err := uc.requestsAnomaliesService.UpsertAnomalies(ctx, uc.responseDtosToEntities(modelAnomalyPredictions)); err != nil {
			z.Errorw(
				"Failed to upsert requests anomalies",
				"error", err.Error(),
			)

			t2 = nil

			break
		}
	}

	return false
}

func (uc *EnrichRequestsWithAnomaliesUsecase) requestEntitiesToDtos(es []*entity.Request) []*dto.GetAnomalyFromModelRequestDto {
	dtos := make([]*dto.GetAnomalyFromModelRequestDto, len(es))

	for i := range es {
		dtos[i] = (*dto.GetAnomalyFromModelRequestDto)(es[i])
	}

	return dtos
}

func (uc *EnrichRequestsWithAnomaliesUsecase) responseDtosToEntities(dtos []*dto.GetAnomalyFromModelResponseDto) []*entity.RequestAnomaly {
	es := make([]*entity.RequestAnomaly, len(dtos))

	for i := range dtos {
		es[i] = (*entity.RequestAnomaly)(dtos[i])
	}

	return es
}
