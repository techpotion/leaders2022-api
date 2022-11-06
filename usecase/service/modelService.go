package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
	"go.uber.org/zap"
)

type ModelService interface {
	GetSinglePrediction(ctx context.Context, req *dto.GetAnomalyFromModelRequestDto) (*dto.GetAnomalyFromModelResponseDto, error)
	GetPredictions(ctx context.Context, req []*dto.GetAnomalyFromModelRequestDto) ([]*dto.GetAnomalyFromModelResponseDto, error)
}

type modelPythonService struct {
	url      string
	endpoint string
}

func NewModelPythonService(url, endpoint string) *modelPythonService {
	return &modelPythonService{url: url, endpoint: endpoint}
}

func (s *modelPythonService) GetSinglePrediction(
	ctx context.Context,
	req *dto.GetAnomalyFromModelRequestDto,
) (*dto.GetAnomalyFromModelResponseDto, error) {
	return nil, fmt.Errorf("not implemented")
}

// nolint // ignore noctx
func (s *modelPythonService) GetPredictions(
	ctx context.Context,
	req []*dto.GetAnomalyFromModelRequestDto,
) ([]*dto.GetAnomalyFromModelResponseDto, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := http.Post(s.url+s.endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to get response from model service :%w", err)
	}
	defer resp.Body.Close()

	predictions := make([]*dto.GetAnomalyFromModelResponseDto, 0, len(req))

	if err := json.NewDecoder(resp.Body).Decode(&predictions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return predictions, nil
}

type modelMockedService struct{}

func NewModelMockedService() *modelMockedService {
	return &modelMockedService{}
}

//nolint // mock
func (s *modelMockedService) GetSinglePrediction(
	ctx context.Context,
	req *dto.GetAnomalyFromModelRequestDto,
) (*dto.GetAnomalyFromModelResponseDto, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("failed to generate number: %w", err)
	}
	ri := int(nBig.Int64())

	if ri != 1 {
		return &dto.GetAnomalyFromModelResponseDto{
			RootID:       req.RootID,
			IsAnomaly:    false,
			AnomalyCases: make([]int, 0, 1),
		}, nil
	}

	return &dto.GetAnomalyFromModelResponseDto{
		RootID:       req.RootID,
		IsAnomaly:    true,
		AnomalyCases: getRandomMockedCases(),
	}, nil
}

func (s *modelMockedService) GetPredictions(
	ctx context.Context,
	req []*dto.GetAnomalyFromModelRequestDto,
) ([]*dto.GetAnomalyFromModelResponseDto, error) {
	resps := make([]*dto.GetAnomalyFromModelResponseDto, 0, len(req))

	for _, r := range req {
		resp, _ := s.GetSinglePrediction(ctx, r)
		resps = append(resps, resp)
	}

	return resps, nil
}

//nolint // mock
func getRandomMockedCases() []int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(12))
	if err != nil {
		zap.S().Errorw("failed to generate number", "error", err.Error())
		return nil
	}
	ri := int(nBig.Int64())

	switch ri {
	case 0:
		return []int{1}
	case 1:
		return []int{1, 2}
	case 2:
		return []int{1, 3}
	case 3:
		return []int{1, 4}
	case 4:
		return []int{1, 5}
	case 5:
		return []int{2}
	case 6:
		return []int{2, 3}
	case 7:
		return []int{2, 4}
	case 8:
		return []int{3, 5}
	case 9:
		return []int{4, 5}
	case 10:
		return []int{5}
	case 11:
		return []int{3, 4, 5}
	default:
		return []int{4}
	}
}
