package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Get rating of anomalies groupped by deffect categories
// @Description returns rating of anomalies groupped by serving companies based on deffect categories
// @Tags        dashboard
// @Accept      json
// @Produce     json
// @Param       region        query    string true "Region"
// @Param       datetime_from query    string true "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to   query    string true "Upper time limit (RFC3339 formatted)"
// @Success     200           {object} dto.CountAnomaliesByDeffectCategoriesResponseDTO
// @Failure     400           {object} dto.CountAnomaliesByDeffectCategoriesResponseDTO
// @Failure     500           {object} dto.CountAnomaliesByDeffectCategoriesResponseDTO
// @Router      /dashboard/anomalies/ratings/deffect_categories       [get]
func (s *Server) CountAnomaliesByDeffectCategoriesHandler(ctx *gin.Context) {
	z := zap.S().With("context", "CountAnomaliesByDeffectCategoriesHandler")

	req := &dto.CountAnomaliesByDeffectCategoriesRequestDTO{}
	resp := &dto.CountAnomaliesByDeffectCategoriesResponseDTO{Success: true}

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

	c, err := s.countAnomaliesByDeffectCategoriesUsecase.Execute(ctx, req.DateTimeFrom, req.DateTimeTo, req.Region)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.DeffectCategories = c
	resp.Count = len(resp.DeffectCategories)

	ctx.JSON(http.StatusOK, resp)
}
