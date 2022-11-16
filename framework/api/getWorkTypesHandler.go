package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get work types
// @Description get a list of work types
// @Tags        filters
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.GetWorkTypesResponseDTO
// @Failure     500 {object} dto.GetWorkTypesResponseDTO
// @Router      /work_types [get]
func (s *Server) GetWorkTypesHandler(ctx *gin.Context) {
	resp := &dto.GetWorkTypesResponseDTO{Success: true}

	workTypes, err := s.getFiltersUsecase.Execute(ctx.Request.Context(), "work_type_done")
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.WorkTypes = workTypes
	resp.Count = len(resp.WorkTypes)

	ctx.JSON(http.StatusOK, resp)
}
