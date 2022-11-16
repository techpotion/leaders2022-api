package api

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/config"
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/logger"
	"gitlab.com/techpotion/leadershack2022/api/usecase"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	engine *gin.Engine

	getFiltersUsecase                        *usecase.GetFiltersUsecase
	getPointsWithFiltersUsecase              *usecase.GetPointsWithFiltersUsecase
	getRequestsByIdsUsecase                  *usecase.GetRequestsByIdsUsecase
	getRequestsAnomaliesByIdsUsecase         *usecase.GetRequestsAnomaliesByIdsUsecase
	getRequestsFullUsecase                   *usecase.GetRequestsFullUsecase
	countPointsWithFiltersUsecase            *usecase.CountPointsWithFiltersUsecase
	countRequestsFullUsecase                 *usecase.CountRequestsFullUsecase
	setCustomRequestAnomalyUsecase           *usecase.SetCustomRequestAnomalyUsecase
	countAnomaliesUsecase                    *usecase.CountAnomaliesUsecase
	countAnomaliesGrouppedUsecase            *usecase.CountAnomaliesGrouppedUsecase
	getAnomaliesAmountDynamicsUsecase        *usecase.GetAnomaliesAmountDynamicsUsecase
	countAnomaliesCountByOwnerCompanies      *usecase.CountAnomaliesByOwnerCompaniesUsecase
	countAnomaliesByServingCompaniesUsecase  *usecase.CountAnomaliesByServingCompaniesUsecase
	countAnomaliesByDeffectCategoriesUsecase *usecase.CountAnomaliesByDeffectCategoriesUsecase
	getEfficiencyPlotUsecase                 *usecase.GetEfficiencyPlotUsecase
	getRegionAreaUsecase                     *usecase.GetRegionAreaUsecase
}

func NewServer(
	cfg *config.Config,
	getFiltersUsecase *usecase.GetFiltersUsecase,
	getPointsWithFiltersUsecase *usecase.GetPointsWithFiltersUsecase,
	getRequestsByIdsUsecase *usecase.GetRequestsByIdsUsecase,
	getRequestsAnomaliesByIdsUsecase *usecase.GetRequestsAnomaliesByIdsUsecase,
	getRequestsFullUsecase *usecase.GetRequestsFullUsecase,
	countPointsWithFiltersUsecase *usecase.CountPointsWithFiltersUsecase,
	countRequestsFullUsecase *usecase.CountRequestsFullUsecase,
	setCustomRequestAnomalyUsecase *usecase.SetCustomRequestAnomalyUsecase,
	countAnomaliesUsecase *usecase.CountAnomaliesUsecase,
	countAnomaliesGrouppedUsecase *usecase.CountAnomaliesGrouppedUsecase,
	getAnomaliesAmountDynamicsUsecase *usecase.GetAnomaliesAmountDynamicsUsecase,
	countAnomaliesCountByOwnerCompanies *usecase.CountAnomaliesByOwnerCompaniesUsecase,
	countAnomaliesByServingCompaniesUsecase *usecase.CountAnomaliesByServingCompaniesUsecase,
	countAnomaliesByDeffectCategoriesUsecase *usecase.CountAnomaliesByDeffectCategoriesUsecase,
	getEfficiencyPlotUsecase *usecase.GetEfficiencyPlotUsecase,
	getRegionAreaUsecase *usecase.GetRegionAreaUsecase,
) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	corsOpt := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	})
	engine.Use(logger.GinRecovery(zap.L(), true), corsOpt)

	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.ServerHost, cfg.ServerPort),
		Handler:           engine,
		ReadTimeout:       cfg.ServerReadTimeout,
		ReadHeaderTimeout: cfg.ServerReadTimeout,
		WriteTimeout:      cfg.ServerWriteTimeout,
	}

	s := &Server{
		server:                                   srv,
		engine:                                   engine,
		getFiltersUsecase:                        getFiltersUsecase,
		getPointsWithFiltersUsecase:              getPointsWithFiltersUsecase,
		getRequestsByIdsUsecase:                  getRequestsByIdsUsecase,
		getRequestsAnomaliesByIdsUsecase:         getRequestsAnomaliesByIdsUsecase,
		getRequestsFullUsecase:                   getRequestsFullUsecase,
		countPointsWithFiltersUsecase:            countPointsWithFiltersUsecase,
		countRequestsFullUsecase:                 countRequestsFullUsecase,
		setCustomRequestAnomalyUsecase:           setCustomRequestAnomalyUsecase,
		countAnomaliesUsecase:                    countAnomaliesUsecase,
		countAnomaliesGrouppedUsecase:            countAnomaliesGrouppedUsecase,
		getAnomaliesAmountDynamicsUsecase:        getAnomaliesAmountDynamicsUsecase,
		countAnomaliesCountByOwnerCompanies:      countAnomaliesCountByOwnerCompanies,
		countAnomaliesByServingCompaniesUsecase:  countAnomaliesByServingCompaniesUsecase,
		countAnomaliesByDeffectCategoriesUsecase: countAnomaliesByDeffectCategoriesUsecase,
		getEfficiencyPlotUsecase:                 getEfficiencyPlotUsecase,
		getRegionAreaUsecase:                     getRegionAreaUsecase,
	}

	s.SetupRouter()

	return s, nil
}

func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			zap.S().Fatalw("Server error", "context", "server", "error", err.Error())
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
