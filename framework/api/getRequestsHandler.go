package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Get requests
// @Description get a list of requests by root_ids
// @Tags        requests
// @Accept      json
// @Produce     json
// @Param       root_ids query    []string false "Root IDs"
// @Success     200      {object} dto.GetRequestsByIdsResponseDto
// @Failure     400      {object} dto.GetRequestsByIdsResponseDto
// @Failure     500      {object} dto.GetRequestsByIdsResponseDto
// @Router      /requests [get]
func (s *Server) GetRequestsHandler(ctx *gin.Context) {
	z := zap.S().With("context", "GetRequestsHandler")

	req := &dto.GetRequestsByIdsRequestDto{}
	resp := &dto.GetRequestsByIdsResponseDto{Success: true}

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

	ids := strings.Split(req.RootIds, ",")

	if len(ids) == 0 {
		resp.Success = false
		resp.Error = "empty ids array provided"
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	requests, err := s.getRequestsByIdsUsecase.Execute(ctx.Request.Context(), ids)
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Requests = requests
	resp.Count = len(resp.Requests)

	ctx.JSON(http.StatusOK, resp)
}
