package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get region area
// @Description get region's area geojson polygon
// @Tags        regions
// @Accept      json
// @Produce     json
// @Param       region path     string true "Region"
// @Success     200    {object} dto.GetRegionAreaResponseDTO
// @Failure     404    {object} dto.GetRegionAreaResponseDTO
// @Failure     500    {object} dto.GetRegionAreaResponseDTO
// @Router      /regions/{region}/area [get]
func (s *Server) GetRegionAreaHandler(ctx *gin.Context) {
	resp := &dto.GetRegionAreaResponseDTO{Success: true}

	region, ok := ctx.Params.Get("region")
	if !ok {
		resp.Error = "region was not found in query"
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	area, err := s.getRegionAreaUsecase.Execute(ctx.Request.Context(), region)
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	if area == "" {
		resp.Error = fmt.Sprintf("region area was not found for region %s", region)
		resp.Success = false
		ctx.JSON(http.StatusNotFound, resp)

		return
	}

	resp.AreaPloygonGeoJSON = area

	ctx.JSON(http.StatusOK, resp)
}
