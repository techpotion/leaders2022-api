package di

import (
	"context"
	"fmt"
	"net"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/config"
	"gitlab.com/techpotion/leadershack2022/api/framework/api"
	"gitlab.com/techpotion/leadershack2022/api/framework/scheduler"
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/database"
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/logger"
	"gitlab.com/techpotion/leadershack2022/api/usecase"
	"gitlab.com/techpotion/leadershack2022/api/usecase/repository"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type DI struct {
	Config *config.Config
	Server *api.Server

	Scheduler *scheduler.Scheduler
	Jobs      []scheduler.Job

	Infrastructure struct {
		Database *database.Postgres
	}
	Repositories struct {
		MoscowRegionsRepository       repository.MoscowRegionsRepository
		HCSRepository                 repository.HCSRepository
		LastAnomalyCheckJobRepository repository.LastAnomalyCheckJobRepository
		RequestsAnomaliesRepository   repository.RequestsAnomaliesRepository
	}
	Services struct {
		MoscowRegionsService       service.MoscowRegionsService
		HCSService                 service.HCSService
		LastAnomalyCheckJobService service.LastAnomalyCheckJobService
		RequestsAnomaliesService   service.RequestsAnomaliesService
		ModelService               service.ModelService
	}
	Usecases struct {
		GetUniqueMoscowRegionsUsecase      *usecase.GetUniqueMoscowRegionsUsecase
		GetPointsWithFiltersUsecase        *usecase.GetPointsWithFiltersUsecase
		GetRequestsByIdsUsecase            *usecase.GetRequestsByIdsUsecase
		GetRequestsAnomaliesByIdsUsecase   *usecase.GetRequestsAnomaliesByIdsUsecase
		EnrichRequestsWithAnomaliesUsecase *usecase.EnrichRequestsWithAnomaliesUsecase
		GetRequestsFullUsecase             *usecase.GetRequestsFullUsecase
		CountPointsWithFiltersUsecase      *usecase.CountPointsWithFiltersUsecase
		CountRequestsFullUsecase           *usecase.CountRequestsFullUsecase
		SetCustomRequestAnomalyUsecase     *usecase.SetCustomRequestAnomalyUsecase
	}
}

func NewDI(ctx context.Context) (*DI, error) {
	di := new(DI)

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config init error: %w", err)
	}

	di.Config = cfg

	err = logger.NewLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("logger init error: %w", err)
	}

	pgDB, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("pg db init error: %w", err)
	}

	di.Infrastructure.Database = pgDB

	di.initRepositories(pgDB)

	di.initServices()

	di.initUsecases()

	server, err := api.NewServer(
		cfg,
		di.Usecases.GetUniqueMoscowRegionsUsecase,
		di.Usecases.GetPointsWithFiltersUsecase,
		di.Usecases.GetRequestsByIdsUsecase,
		di.Usecases.GetRequestsAnomaliesByIdsUsecase,
		di.Usecases.GetRequestsFullUsecase,
		di.Usecases.CountPointsWithFiltersUsecase,
		di.Usecases.CountRequestsFullUsecase,
		di.Usecases.SetCustomRequestAnomalyUsecase,
	)
	if err != nil {
		return nil, fmt.Errorf("server init error: %w", err)
	}

	di.Jobs = append(
		di.Jobs,
		scheduler.NewAnomaliesDailyJob("30 1 * * *", di.Usecases.EnrichRequestsWithAnomaliesUsecase),
	)
	di.Scheduler = scheduler.NewScheduler(di.Jobs...)

	di.Server = server

	// nolint:gocritic // test cases
	// reqs, _ := di.Usecases.GetRequestsByIdsUsecase.Execute(ctx, []string{"73417411", "73417412"})
	//
	// dtoReqs := make([]*dto.GetAnomalyFromModelRequestDto, 0, len(reqs))
	// for _, r := range reqs {
	// 	t := time.Now()
	// 	dtoReqs = append(dtoReqs, &dto.GetAnomalyFromModelRequestDto{Request: r, DateOfPreviousRequestClose: &t})
	// }
	//
	// resp, err := di.Services.ModelService.GetPredictions(ctx, dtoReqs)
	// if err != nil {
	// 	log.Panic(err)
	// }
	//
	// for _, aa := range resp {
	// 	fmt.Printf("%+v\n", aa)
	// 	b, _ := json.Marshal(aa)
	// 	fmt.Println(string(b))
	// }
	//
	// t1, _ := time.Parse("2006-01-02 15:04:05", "2021-01-14 15:50:27")
	// t2, _ := time.Parse("2006-01-02 15:04:05", "2021-01-15 15:50:27")
	//
	// fmt.Println(di.Services.HCSService.CountRequestsByClosureTime(ctx, t1, t2))
	// fmt.Println(di.Services.HCSService.GetRequestsByClosureTime(ctx, t1, t2, 5, 0))
	//
	di.Usecases.EnrichRequestsWithAnomaliesUsecase.Execute(ctx)

	return di, nil
}

func (di *DI) Start(ctx context.Context) context.Context {
	z := zap.S().With("context", "di.Start")

	z.Info("Starting scheduler")

	di.Scheduler.Start()

	z.Infow("Starting HTTP server", "URI", net.JoinHostPort(di.Config.ServerHost, di.Config.ServerPort))

	di.Server.Start()

	return ctx
}

func (di *DI) Stop() {
	z := zap.S().With("context", "di.Stop")

	z.Info("Closing HTTP server")

	//nolint:gomnd // 15 seconds
	sctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	if err := di.Server.Stop(sctx); err != nil {
		z.Errorw("Error closing HTTP server", "error", err)
	}

	cancel()

	z.Info("Closing scheduler")

	sdlrCtx := di.Scheduler.Stop()
	<-sdlrCtx.Done()

	z.Info("Closing postgres connection")

	di.Infrastructure.Database.Pool.Close()
}
