package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Get rating of anomalies groupped by owner companies
// @Description returns rating of anomalies groupped by owner companies based on anomalies count
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       region        query    string true "Region"
// @Param       datetime_from query    string true "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to   query    string true "Upper time limit (RFC3339 formatted)"
// @Success     200           {object} dto.CountAnomaliesByOwnerCompaniesResponseDTO
// @Failure     400           {object} dto.CountAnomaliesByOwnerCompaniesResponseDTO
// @Failure     500           {object} dto.CountAnomaliesByOwnerCompaniesResponseDTO
// @Router      /dashboard/anomalies/ratings/owner_companies       [get]
func (s *Server) CountAnomaliesByOwnerCompaniesHandler(ctx *gin.Context) {
	z := zap.S().With("context", "CountAnomaliesByOwnerCompaniesHandler")

	req := &dto.CountAnomaliesByOwnerCompaniesRequestDTO{}
	resp := &dto.CountAnomaliesByOwnerCompaniesResponseDTO{Success: true}

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

	c, err := s.countAnomaliesCountByOwnerCompanies.Execute(ctx, req.DateTimeFrom, req.DateTimeTo, req.Region)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.OwnerCompanies = c
	resp.Count = len(resp.OwnerCompanies)

	ctx.JSON(http.StatusOK, resp)
}
