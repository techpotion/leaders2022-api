package di

import (
	"gitlab.com/techpotion/leadershack2022/api/usecase"
)

func (di *DI) initUsecases() {
	// endpoint usecases
	di.Usecases.GetUniqueMoscowRegionsUsecase = usecase.NewGetUniqueMoscowRegionsUsecase(di.Services.MoscowRegionsService)
	di.Usecases.GetPointsWithFiltersUsecase = usecase.NewGetPointsWithFiltersUsecase(di.Services.HCSService)
	di.Usecases.GetRequestsByIdsUsecase = usecase.NewGetRequestsByIdsUsecase(di.Services.HCSService)
	di.Usecases.GetRequestsAnomaliesByIdsUsecase = usecase.NewGetRequestsAnomaliesByIdsUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.GetRequestsFullUsecase = usecase.NewGetRequestsFullUsecase(di.Services.HCSService)
	di.Usecases.CountPointsWithFiltersUsecase = usecase.NewCountPointsWithFiltersUsecase(di.Services.HCSService)
	di.Usecases.CountRequestsFullUsecase = usecase.NewCountRequestsFullUsecase(di.Services.HCSService)
	di.Usecases.SetCustomRequestAnomalyUsecase = usecase.NewSetCustomRequestAnomalyUsecase(di.Services.RequestsAnomaliesService)

	// internal usecases
	di.Usecases.EnrichRequestsWithAnomaliesUsecase = usecase.NewEnrichRequestsWithAnomaliesUsecase(
		di.Services.LastAnomalyCheckJobService,
		di.Services.HCSService,
		di.Services.ModelService,
		di.Services.RequestsAnomaliesService,
	)
}
