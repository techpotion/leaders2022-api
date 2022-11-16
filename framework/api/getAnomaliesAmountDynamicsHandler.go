package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Get anomalies amount dynamics
// @Description returns anomalies amount dynamics
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       region        query    string true "Region"
// @Param       datetime_from query    string true "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to   query    string true "Upper time limit (RFC3339 formatted)"
// @Success     200           {object} dto.GetAnomaliesAmountDynamicsResponseDTO
// @Failure     400           {object} dto.GetAnomaliesAmountDynamicsResponseDTO
// @Failure     500           {object} dto.GetAnomaliesAmountDynamicsResponseDTO
// @Router      /dashboard/anomalies/amount_dynamics       [get]
func (s *Server) GetAnomaliesAmountDynamicsHandler(ctx *gin.Context) {
	z := zap.S().With("context", "GetAnomaliesAmountDynamicsHandler")

	req := &dto.GetAnomaliesAmountDynamicsRequestDTO{}
	resp := &dto.GetAnomaliesAmountDynamicsResponseDTO{Success: true}

	if err := ctx.ShouldBind(req); err != nil {
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

	d, err := s.getAnomaliesAmountDynamicsUsecase.Execute(ctx, req.DateTimeFrom, req.DateTimeTo, req.Region)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Dynamics = d
	resp.Count = len(resp.Dynamics)

	ctx.JSON(http.StatusOK, resp)
}
