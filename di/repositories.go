package di

import (
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/database"
	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
)

func (di *DI) initRepositories(pgDB *database.Postgres) {
	di.Repositories.HCSRepository = repository.NewHCSPostgresRepository(pgDB)
	di.Repositories.LastAnomalyCheckJobRepository = repository.NewLastAnomalyCheckJobPostgresRepository(pgDB)
	di.Repositories.RequestsAnomaliesRepository = repository.NewRequestsAnomaliesPostgresRepository(pgDB)
}
