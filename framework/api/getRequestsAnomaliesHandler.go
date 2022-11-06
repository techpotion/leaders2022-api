package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

func splitToInts(s, sep string) []int {
	tmp := strings.Split(s, sep)
	values := make([]int, 0, len(tmp))

	for _, raw := range tmp {
		v, err := strconv.Atoi(raw)
		if err != nil {
			log.Print(err)
			continue
		}

		values = append(values, v)
	}

	return values
}

// @Summary     Get anomalies
// @Description get a list of requests' anomalies by root_ids
// @Tags        anomalies
// @Accept      json
// @Produce     json
// @Param       root_ids      query    []string true  "Root IDs"
// @Param       anomaly_cases query    []int    false "Anomaly Cases"
// @Success     200           {object} dto.GetRequestsAnomaliesByIdsResponseDto
// @Failure     400           {object} dto.GetRequestsAnomaliesByIdsResponseDto
// @Failure     500           {object} dto.GetRequestsAnomaliesByIdsResponseDto
// @Router      /anomalies [get]
func (s *Server) GetRequestsAnomaliesHandler(ctx *gin.Context) {
	z := zap.S().With("context", "GetRequestsHandler")

	req := &dto.GetRequestsAnomaliesByIdsRequestDto{}
	resp := &dto.GetRequestsAnomaliesByIdsResponseDto{Success: true}

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

	cases := []int{}
	if req.AnomalyCases != "" {
		cases = splitToInts(req.AnomalyCases, ",")
	}

	if len(ids) == 0 {
		resp.Success = false
		resp.Error = "empty root ids array provided"
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	anomalies, err := s.getRequestsAnomaliesByIdsUsecase.Execute(ctx.Request.Context(), ids, cases)
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Anomalies = anomalies
	resp.Count = len(resp.Anomalies)

	ctx.JSON(http.StatusOK, resp)
}
