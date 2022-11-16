package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
)

// @Summary     Get efficiency plot for dispatcher
// @Description get an efficiency plot for dispatcher
// @Tags        dashboard
// @Param       dispatcher    query    string true "Dispatcher (ОДС)"
// @Param       datetime_from query    string true "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to   query    string true "Upper time limit (RFC3339 formatted)"
// @Success     200           {object} dto.GetEfficiencyPlotResponseDTO
// @Failure     500           {object} dto.GetEfficiencyPlotResponseDTO
// @Router      /dashboard/plots/efficiency       [get]
func (s *Server) GetEfficiencyPlotHandler(ctx *gin.Context) {
	req := &dto.GetEfficiencyPlotRequestDTO{}
	resp := &dto.GetEfficiencyPlotResponseDTO{Success: true}

	if err := ctx.ShouldBind(req); err != nil {
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

	scs, err := s.getEfficiencyPlotUsecase.Execute(ctx.Request.Context(), req)
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.Filename = &scs

	ctx.JSON(http.StatusOK, resp)
}
