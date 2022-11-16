package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get regions
// @Description get a list of moscow regions
// @Tags        filters
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.GetRegionsResponseDTO
// @Failure     500 {object} dto.GetRegionsResponseDTO
// @Router      /regions [get]
func (s *Server) GetRegionsHandler(ctx *gin.Context) {
	resp := &dto.GetRegionsResponseDTO{Success: true}

	regions, err := s.getFiltersUsecase.Execute(ctx.Request.Context(), "hood")
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.Regions = regions
	resp.Count = len(resp.Regions)

	ctx.JSON(http.StatusOK, resp)
}
