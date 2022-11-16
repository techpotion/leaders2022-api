package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     DEPRECATED Count requests fulls by selected area and region
// @Description returns amount of requests full
// @Tags        requests_full
// @Accept      json
// @Produce     json
// @Param       region           query    string  true  "Region"
// @Param       datetime_from    query    string  true  "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to      query    string  true  "Upper time limit (RFC3339 formatted)"
// @Param       x_min            query    float32 false "Bottom left X of screen"
// @Param       y_min            query    float32 false "Bottom left Y of screen"
// @Param       x_max            query    float32 false "Top right X of screen"
// @Param       y_max            query    float32 false "Top right Y of screen"
// @Param       urgency_category query    string  false "Urgency Category"
// @Param       anomaly_cases    query    []int   false "Anomaly Cases"
// @Success     200              {object} dto.CountRequestsFullResponseDTO
// @Failure     400              {object} dto.CountRequestsFullResponseDTO
// @Failure     500              {object} dto.CountRequestsFullResponseDTO
// @Router      /requests_full/count        [get]
func (s *Server) CountRequestsFullHandler(ctx *gin.Context) {
	z := zap.S().With("context", "CountRequestsFullHandler")

	req := &dto.CountRequestsFullRequestDTO{}
	resp := &dto.CountRequestsFullResponseDTO{Success: true}

	if err := ctx.ShouldBind(req); err != nil {
		z.Errorw("failed to bind query", "error", err.Error())

		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	if req.XMin == nil ||
		req.Ymin == nil ||
		req.XMax == nil ||
		req.YMax == nil {
		req.XMin, req.Ymin, req.XMax, req.YMax = &minXMocked, &minYMocked, &maxXMocked, &maxYMocked
	}

	if err := validation.Struct(req); err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	count, err := s.countRequestsFullUsecase.Execute(ctx, req)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Count = count

	ctx.JSON(http.StatusOK, resp)
}
