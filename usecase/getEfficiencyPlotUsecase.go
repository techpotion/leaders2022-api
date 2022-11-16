package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"gitlab.com/techpotion/leadershack2022/api/usecase/service"
	"go.uber.org/zap"
)

type GetEfficiencyPlotUsecase struct {
	hcsService  service.HCSService
	fileService service.FileService
	plotService service.PlotService
}

func NewGetEfficiencyPlotUsecase(
	hcsService service.HCSService,
	fileService service.FileService,
	plotService service.PlotService,
) *GetEfficiencyPlotUsecase {
	return &GetEfficiencyPlotUsecase{hcsService: hcsService, fileService: fileService, plotService: plotService}
}

func (uc *GetEfficiencyPlotUsecase) Execute(ctx context.Context, req *dto.GetEfficiencyPlotRequestDTO) (string, error) {
	z := zap.S().With("context", "GetEfficiencyPlotUsecase")

	if req == nil {
		return "", fmt.Errorf("empty request")
	}

	requests, err := uc.hcsService.GetRequestsByDispatcher(ctx, req.DateTimeFrom, req.DateTimeTo, req.Dispatcher)
	if err != nil {
		z.Errorw("failed to get requests by dispatcher", "error", err.Error())
		return "", fmt.Errorf("failed to get requests by dispatcher")
	}

	dtos, err := uc.requestsToDto(requests)
	if err != nil {
		z.Errorw("failed to convert requests to dtos", "error", err.Error())
		return "", fmt.Errorf("failed to convert requests to dtos")
	}

	data, err := uc.plotService.GetPlotSvgString(ctx, dtos)
	if err != nil {
		z.Errorw("failed to get plots from service", "error", err.Error())
		return "", fmt.Errorf("failed to get plots from plot service")
	} else if data == nil {
		z.Errorw("failed to get plots from service: empty path received")
		return "", fmt.Errorf("failed to get plots from service: empty path received")
	}

	filename := uuid.New().String() + ".svg"
	if err := uc.fileService.WriteStringToFile(data, filename); err != nil {
		z.Errorw("failed to write plot to file", "error", err.Error())
		return "", fmt.Errorf("failed to write plot to file: %w", err)
	}

	return filename, nil
}

func (uc *GetEfficiencyPlotUsecase) requestsToDto(requests []*entity.Request) ([]*dto.GetPlotRequestDTO, error) {
	dtos := make([]*dto.GetPlotRequestDTO, len(requests))

	for i := range requests {
		if requests[i].DispetchersNumber == nil {
			return nil, fmt.Errorf("empty dispatcher in array")
		}

		dtos[i] = &dto.GetPlotRequestDTO{
			RootID:            requests[i].RootID,
			DispatchersNumber: *requests[i].DispetchersNumber,
			Efficiency:        requests[i].Effeciency,
			ClosureDate:       requests[i].ClosureDate,
			DateOfReview:      requests[i].DateOfReview,
			GradeForService:   requests[i].GradeForService,
		}
	}

	return dtos, nil
}
