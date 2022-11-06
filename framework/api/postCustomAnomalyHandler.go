package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Set custom anomaly marked by user
// @Description sets anomaly data from provided in request body json
// @Tags        anomalies
// @Accept      json
// @Produce     json
// @Param       data body     dto.PostCustomAnomaliesRequestDTO true "Custom anomaly body"
// @Success     201  {object} dto.CountPointsResponseDTO
// @Failure     400  {object} dto.CountPointsResponseDTO
// @Failure     500  {object} dto.CountPointsResponseDTO
// @Router      /custom_requests_anomalies        [post]
func (s *Server) PostCustomAnomalyHandler(ctx *gin.Context) {
	z := zap.S().With("context", "CountPointsHandler")

	req := &dto.PostCustomAnomaliesRequestDTO{}
	resp := &dto.PostCustomAnomaliesResponseDTO{Success: true}

	if err := ctx.ShouldBindJSON(req); err != nil {
		z.Errorw("failed to bind query", "error", err.Error())

		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	if err := validation.Struct(req); err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	if err := s.setCustomRequestAnomalyUsecase.Execute(ctx, req.RootID, *req.IsAnomaly); err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
