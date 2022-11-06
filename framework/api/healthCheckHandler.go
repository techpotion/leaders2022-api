package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary     Health check
// @Description simple endpoint to check whether service is alive
// @Tags        health
// @Router      /health [get]
func (s *Server) HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
