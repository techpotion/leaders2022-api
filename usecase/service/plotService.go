package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

type PlotService interface {
	GetPlotSvgString(ctx context.Context, req []*dto.GetPlotRequestDTO) ([]byte, error)
}

type plotPythonService struct {
	url      string
	endpoint string
}

func NewPlotPythonService(url, endpoint string) *plotPythonService {
	return &plotPythonService{url: url, endpoint: endpoint}
}

// nolint // @TODO fix
func (s *plotPythonService) GetPlotSvgString(ctx context.Context, req []*dto.GetPlotRequestDTO) ([]byte, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	resp, err := http.Post(s.url+s.endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to get response from plot service :%w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get response from microservice: status code: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return b, nil
}
