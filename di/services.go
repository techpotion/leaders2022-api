package di

import "gitlab.com/techpotion/leadershack2022/api/usecase/service"

func (di *DI) initServices() {
	di.Services.MoscowRegionsService = service.NewMoscowRegionsService(di.Repositories.MoscowRegionsRepository)
	di.Services.HCSService = service.NewHCSService(di.Repositories.HCSRepository)
	di.Services.LastAnomalyCheckJobService = service.NewLastAnomalyCheckJobService(di.Repositories.LastAnomalyCheckJobRepository)
	di.Services.RequestsAnomaliesService = service.NewRequestsAnomaliesService(di.Repositories.RequestsAnomaliesRepository)
	di.Services.ModelService = service.NewModelPythonService(di.Config.ModelMicroserviceURI, di.Config.ModelMicroservicePredictEndpoint)
}
