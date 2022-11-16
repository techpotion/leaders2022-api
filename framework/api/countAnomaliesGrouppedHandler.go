package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Count and get amount anomalies groupped by cases
// @Description returns amount of anomalies groupped by cases
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       region        query    string true "Region"
// @Param       datetime_from query    string true "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to   query    string true "Upper time limit (RFC3339 formatted)"
// @Success     200           {object} dto.CountAnomaliesGrouppedResponseDTO
// @Failure     400           {object} dto.CountAnomaliesGrouppedResponseDTO
// @Failure     500           {object} dto.CountAnomaliesGrouppedResponseDTO
// @Router      /dashboard/anomalies/count_groupped       [get]
func (s *Server) CountAnomaliesGrouppedHandler(ctx *gin.Context) {
	z := zap.S().With("context", "CountAnomaliesGrouppedHandler")

	req := &dto.CountAnomaliesGrouppedRequestDTO{}
	resp := &dto.CountAnomaliesGrouppedResponseDTO{Success: true}

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

	m, err := s.countAnomaliesGrouppedUsecase.Execute(ctx, req.DateTimeFrom, req.DateTimeTo, req.Region)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.GroupsCounts = m

	ctx.JSON(http.StatusOK, resp)
}
