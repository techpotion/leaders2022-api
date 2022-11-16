package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get owner companies
// @Description get a list of owner companies
// @Tags        filters
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.GetOwnerCompaniesResponseDTO
// @Failure     500 {object} dto.GetOwnerCompaniesResponseDTO
// @Router      /owner_companies [get]
func (s *Server) GetOwnerCompanies(ctx *gin.Context) {
	resp := &dto.GetOwnerCompaniesResponseDTO{Success: true}

	scs, err := s.getFiltersUsecase.Execute(ctx.Request.Context(), "owner_company")
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.OwnerCompanies = scs
	resp.Count = len(resp.OwnerCompanies)

	ctx.JSON(http.StatusOK, resp)
}
