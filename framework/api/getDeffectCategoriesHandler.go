package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get deffect category_names
// @Description get a list of deffect category_names
// @Tags        filters
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.GetDeffectCategoriesResponseDTO
// @Failure     500 {object} dto.GetDeffectCategoriesResponseDTO
// @Router      /deffect_categories [get]
func (s *Server) GetDeffectCategories(ctx *gin.Context) {
	resp := &dto.GetDeffectCategoriesResponseDTO{Success: true}

	scs, err := s.getFiltersUsecase.Execute(ctx.Request.Context(), "deffect_category_name")
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.DeffectCategories = scs
	resp.Count = len(resp.DeffectCategories)

	ctx.JSON(http.StatusOK, resp)
}
