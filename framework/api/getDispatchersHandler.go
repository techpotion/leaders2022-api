package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

// @Summary     Get dispatchers
// @Description get a list of dispatchers
// @Tags        filters
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.GetDispatchersResponseDTO
// @Failure     500 {object} dto.GetDispatchersResponseDTO
// @Router      /dispatchers [get]
func (s *Server) GetDispatchersHandler(ctx *gin.Context) {
	resp := &dto.GetDispatchersResponseDTO{Success: true}

	scs, err := s.getFiltersUsecase.Execute(ctx.Request.Context(), "dispetchers_number")
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	resp.Dispatchers = scs
	resp.Count = len(resp.Dispatchers)

	ctx.JSON(http.StatusOK, resp)
}
