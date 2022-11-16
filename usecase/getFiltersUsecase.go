package usecase

import (
	"context"
	"fmt"

	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetFiltersUsecase struct {
	hcsService   service.HCSService
	responsesMap map[string]func(context.Context) ([]string, error)
}

func NewGetFiltersUsecase(hcsService service.HCSService) *GetFiltersUsecase {
	uc := &GetFiltersUsecase{hcsService: hcsService}

	uc.responsesMap = make(map[string]func(context.Context) ([]string, error))
	uc.responsesMap["hood"] = uc.hcsService.GetUniqueRegions
	uc.responsesMap["serving_company"] = uc.hcsService.GetUniqueServingCompanies
	uc.responsesMap["owner_company"] = uc.hcsService.GetUniqueOwnerCompanies
	uc.responsesMap["deffect_category_name"] = uc.hcsService.GetUniqueDeffectCategories
	uc.responsesMap["work_type_done"] = uc.hcsService.GetUniqueWorkTypes
	uc.responsesMap["dispetchers_number"] = uc.hcsService.GetUniqueDispatchers

	return uc
}

func (uc *GetFiltersUsecase) Execute(ctx context.Context, column string) ([]string, error) {
	z := zap.S().With("context", "GetUniqueMoscowRegionsUsecase")

	f, ok := uc.responsesMap[column]
	if !ok {
		return nil, fmt.Errorf("failed to get find service method by column name: %s", column)
	}

	resp, err := f(ctx)
	if err != nil {
		z.Errorw("failed to get unique column values", "error", err.Error(), "column", column)
	}

	return resp, err
}
