package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Get points by selected area and region
// @Description returns a list of points in selected area
// @Tags        requests_full
// @Accept      json
// @Produce     json
// @Param       region           query    string  true  "Region"
// @Param       datetime_from    query    string  true  "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to      query    string  true  "Upper time limit (RFC3339 formatted)"
// @Param       x_min            query    float32 true  "Bottom left X of screen"
// @Param       y_min            query    float32 true  "Bottom left Y of screen"
// @Param       x_max            query    float32 true  "Top right X of screen"
// @Param       y_max            query    float32 true  "Top right Y of screen"
// @Param       limit            query    int     true  "Limit"
// @Param       offset           query    int     true  "Offset"
// @Param       urgency_category query    string  false "Urgency Category"
// @Param       anomaly_cases    query    []int   false "Anomaly Cases"
// @Success     200              {object} dto.GetRequestsFullResponseDTO
// @Failure     400              {object} dto.GetRequestsFullResponseDTO
// @Failure     500              {object} dto.GetRequestsFullResponseDTO
// @Router      /requests_full        [get]
func (s *Server) GetRequestsFullHandler(ctx *gin.Context) {
	z := zap.S().With("context", "GetRequestsFullHandler")

	req := &dto.GetRequestsFullRequestDTO{}
	resp := &dto.GetRequestsFullResponseDTO{Success: true}

	if casesStr, ok := ctx.GetQuery("anomaly_cases"); ok {
		req.AnomalyCases = splitToInts(casesStr, ",")
	}

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

	requests, err := s.getRequestsFullUsecase.Execute(ctx, req)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Requests = requests
	resp.Count = len(resp.Requests)

	ctx.JSON(http.StatusOK, resp)
}
