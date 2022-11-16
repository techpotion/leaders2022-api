package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/validation"
	"go.uber.org/zap"
)

// @Summary     Count points by selected area and region
// @Description returns amount of entities
// @Tags        points
// @Accept      json
// @Produce     json
// @Param       region            query    string   true  "Region"
// @Param       datetime_from     query    string   true  "Lower time limit (RFC3339 formatted)"
// @Param       datetime_to       query    string   true  "Upper time limit (RFC3339 formatted)"
// @Param       x_min             query    float32  false "Bottom left X of screen"
// @Param       y_min             query    float32  false "Bottom left Y of screen"
// @Param       x_max             query    float32  false "Top right X of screen"
// @Param       y_max             query    float32  false "Top right Y of screen"
// @Param       serving_company   query    []string false "Serving companies (компания исполнитель)"
// @Param       efficiency        query    []string false "Efficiency (результативность)"
// @Param       grade_for_service query    []string false "Grade for service (оценка качества)"
// @Param       urgency_category  query    []string false "Urgency category (срочность)"
// @Param       work_type         query    []string false "Work type (вид работы)"
// @Param       deffect_category  query    []string false "Deffect category (категория дефекта)"
// @Param       owner_company     query    []string false "Owner company (упр. организация)"
// @Success     200               {object} dto.CountPointsResponseDTO
// @Failure     400               {object} dto.CountPointsResponseDTO
// @Failure     500               {object} dto.CountPointsResponseDTO
// @Router      /points/count        [get]
func (s *Server) CountPointsHandler(ctx *gin.Context) {
	z := zap.S().With("context", "CountPointsHandler")

	req := &dto.CountPointsRequestDTO{}
	resp := &dto.CountPointsResponseDTO{Success: true}

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

	if len(req.ServingCompany) == 1 {
		req.ServingCompany = strings.Split(req.ServingCompany[0], ",")
	}

	if len(req.Efficiency) == 1 {
		req.Efficiency = strings.Split(req.Efficiency[0], ",")
	}

	if len(req.GradeForService) == 1 {
		req.GradeForService = strings.Split(req.GradeForService[0], ",")
	}

	if len(req.UrgencyCategory) == 1 {
		req.UrgencyCategory = strings.Split(req.UrgencyCategory[0], ",")
	}

	if len(req.WorkType) == 1 {
		req.WorkType = strings.Split(req.WorkType[0], ",")
	}

	if len(req.DeffectCategory) == 1 {
		req.DeffectCategory = strings.Split(req.DeffectCategory[0], ",")
	}

	if len(req.OwnerCompany) == 1 {
		req.OwnerCompany = strings.Split(req.OwnerCompany[0], ",")
	}

	count, err := s.countPointsWithFiltersUsecase.Execute(ctx, req)
	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Count = count

	ctx.JSON(http.StatusOK, resp)
}
