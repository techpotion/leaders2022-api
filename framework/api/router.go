package api

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "gitlab.com/techpotion/leadershack2022/api/docs"
)

func (s *Server) SetupRouter() {
	s.engine.GET("/health", s.HealthCheckHandler)

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api/v1"

	apiGroup := s.engine.Group("/api/v1")
	apiGroup.GET("health", s.HealthCheckHandler)
	apiGroup.GET("regions", s.GetRegionsHandler)
	apiGroup.GET("regions/:region/area", s.GetRegionAreaHandler)
	apiGroup.GET("serving_companies", s.GetServingCompanies)
	apiGroup.GET("owner_companies", s.GetOwnerCompanies)
	apiGroup.GET("deffect_categories", s.GetDeffectCategories)
	apiGroup.GET("work_types", s.GetWorkTypesHandler)
	apiGroup.GET("dispatchers", s.GetDispatchersHandler)
	apiGroup.GET("points", s.GetPointsHandler)
	apiGroup.GET("points/count", s.CountPointsHandler)
	apiGroup.GET("requests", s.GetRequestsHandler)
	apiGroup.GET("requests_full", s.GetRequestsFullHandler)
	apiGroup.GET("requests_full/count", s.CountRequestsFullHandler)
	apiGroup.GET("anomalies", s.GetRequestsAnomaliesHandler)
	apiGroup.POST("custom_requests_anomalies", s.PostCustomAnomalyHandler)

	apiGroup.GET("/dashboard/anomalies/count", s.CountAnomaliesHandler)
	apiGroup.GET("/dashboard/anomalies/count_groupped", s.CountAnomaliesGrouppedHandler)
	apiGroup.GET("/dashboard/anomalies/amount_dynamics", s.GetAnomaliesAmountDynamicsHandler)
	apiGroup.GET("/dashboard/anomalies/ratings/deffect_categories", s.CountAnomaliesByDeffectCategoriesHandler)
	apiGroup.GET("/dashboard/anomalies/ratings/serving_companies", s.CountAnomaliesByServingCompaniesHandler)
	apiGroup.GET("/dashboard/anomalies/ratings/owner_companies", s.CountAnomaliesByOwnerCompaniesHandler)
	apiGroup.GET("/dashboard/plots/efficiency", s.GetEfficiencyPlotHandler)

	staticmw := func() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ctx.Header("Content-Disposition", "attachment")
			ctx.Next()
		}
	}

	docs.SwaggerInfo.Host = ""
	if env := os.Getenv("ENV"); env == "test" {
		docs.SwaggerInfo.Host = ""
	}

	staticApiGroup := s.engine.Group("/")
	staticApiGroup.Use(staticmw())
	staticApiGroup.Static("data/", "data")
}
