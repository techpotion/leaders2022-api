package api

import (
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
	apiGroup.GET("points", s.GetPointsHandler)
	apiGroup.GET("points/count", s.CountPointsHandler)
	apiGroup.GET("requests", s.GetRequestsHandler)
	apiGroup.GET("requests_full", s.GetRequestsFullHandler)
	apiGroup.GET("requests_full/count", s.CountRequestsFullHandler)
	apiGroup.GET("anomalies", s.GetRequestsAnomaliesHandler)
	apiGroup.POST("custom_requests_anomalies", s.PostCustomAnomalyHandler)
}
