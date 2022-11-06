package dto

type PostCustomAnomaliesRequestDTO struct {
	RootID    string `form:"root_id"    json:"root_id"    validate:"required"`
	IsAnomaly *bool  `form:"is_anomaly" json:"is_anomaly" validate:"required"`
}

type PostCustomAnomaliesResponseDTO struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
