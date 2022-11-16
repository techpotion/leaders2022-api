package di

import (
	"gitlab.com/techpotion/leadershack2022/api/usecase"
)

func (di *DI) initUsecases() {
	// endpoint usecases
	di.Usecases.GetFiltersUsecase = usecase.NewGetFiltersUsecase(di.Services.HCSService)
	di.Usecases.GetPointsWithFiltersUsecase = usecase.NewGetPointsWithFiltersUsecase(di.Services.HCSService)
	di.Usecases.GetRequestsByIdsUsecase = usecase.NewGetRequestsByIdsUsecase(di.Services.HCSService)
	di.Usecases.GetRequestsAnomaliesByIdsUsecase = usecase.NewGetRequestsAnomaliesByIdsUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.GetRequestsFullUsecase = usecase.NewGetRequestsFullUsecase(di.Services.HCSService)
	di.Usecases.CountPointsWithFiltersUsecase = usecase.NewCountPointsWithFiltersUsecase(di.Services.HCSService)
	di.Usecases.CountRequestsFullUsecase = usecase.NewCountRequestsFullUsecase(di.Services.HCSService)
	di.Usecases.SetCustomRequestAnomalyUsecase = usecase.NewSetCustomRequestAnomalyUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.CountAnomaliesUsecase = usecase.NewCountAnomaliesUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.CountAnomaliesGrouppedUsecase = usecase.NewCountAnomaliesGrouppedUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.GetAnomaliesAmountDynamicsUsecase = usecase.
		NewGetAnomaliesAmountDynamicsUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.CountAnomaliesCountByOwnerCompanies = usecase.
		NewCountAnomaliesByOwnerCompaniesUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.CountAnomaliesByServingCompaniesUsecase = usecase.
		NewCountAnomaliesByServingCompaniesUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.CountAnomaliesByDeffectCategories = usecase.
		NewCountAnomaliesByDeffectCategoriesUsecase(di.Services.RequestsAnomaliesService)
	di.Usecases.GetEfficiencyPlotUsecase = usecase.
		NewGetEfficiencyPlotUsecase(di.Services.HCSService, di.Services.FileService, di.Services.PlotService)
	di.Usecases.GetRegionAreaUsecase = usecase.NewGetRegionAreaUsecase(di.Services.HCSService)

	// internal usecases
	di.Usecases.EnrichRequestsWithAnomaliesUsecase = usecase.NewEnrichRequestsWithAnomaliesUsecase(
		di.Services.LastAnomalyCheckJobService,
		di.Services.HCSService,
		di.Services.ModelService,
		di.Services.RequestsAnomaliesService,
	)
}
