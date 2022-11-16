package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
)

// @Summary     Count and get percent of anomalies
// @Description returns amount and percent of anomalies
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       region        query    string true "Region"
// @Param       datetime_from query    string true "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to   query    string true "Upper time limit (RFC3339 formatted)"
// @Success     200           {object} dto.CountAnomaliesResponseDTO
// @Failure     400           {object} dto.CountAnomaliesResponseDTO
// @Failure     500           {object} dto.CountAnomaliesResponseDTO
// @Router      /dashboard/anomalies/count       [get]
func (s *Server) CountAnomaliesHandler(ctx *gin.Context) {
	req := &dto.CountAnomaliesRequestDTO{}
	resp := &dto.CountAnomaliesResponseDTO{Success: true}

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

	count, perc, err := s.countAnomaliesUsecase.Execute(ctx, req.DateTimeFrom, req.DateTimeTo, req.Region)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	duration := req.DateTimeTo.Sub(req.DateTimeFrom)
	dtfPrev, dttPrev := req.DateTimeFrom.Add(-duration), req.DateTimeFrom

	countPrev, percPrev, err := s.countAnomaliesUsecase.Execute(ctx, dtfPrev, dttPrev, req.Region)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Current = &dto.CountAnomaliesDTO{Count: count, Percent: perc}
	resp.Previous = &dto.CountAnomaliesDTO{Count: countPrev, Percent: percPrev}

	ctx.JSON(http.StatusOK, resp)
}
