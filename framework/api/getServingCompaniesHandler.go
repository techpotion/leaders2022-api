package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get serving companies
// @Description get a list of serving companies
// @Tags        filters
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.GetServingCompaniesResponseDTO
// @Failure     500 {object} dto.GetServingCompaniesResponseDTO
// @Router      /serving_companies [get]
func (s *Server) GetServingCompanies(ctx *gin.Context) {
	resp := &dto.GetServingCompaniesResponseDTO{Success: true}

	scs, err := s.getFiltersUsecase.Execute(ctx.Request.Context(), "serving_company")
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.ServingCompanies = scs
	resp.Count = len(resp.ServingCompanies)

	ctx.JSON(http.StatusOK, resp)
}
