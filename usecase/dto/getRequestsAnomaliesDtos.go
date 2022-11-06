package dto

import "gitlab.com/techpotion/leadershack2022/api/entity"

type GetRequestsAnomaliesByIdsRequestDto struct {
	RootIds      string `form:"root_ids" validate:"required"`
	AnomalyCases string `form:"anomaly_cases"`
}

type GetRequestsAnomaliesByIdsResponseDto struct {
	Anomalies []entity.RequestAnomaly `json:"requests_anomalies"`
	Count     int                     `json:"count"`
	Success   bool                    `json:"success"`
	Error     string                  `json:"error,omitempty"`
}
