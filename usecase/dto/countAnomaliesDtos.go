package dto

import "time"

type CountAnomaliesRequestDTO struct {
	Region       string    `form:"region"          validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type CountAnomaliesDTO struct {
	Count   int     `json:"count"`
	Percent float32 `json:"percent"`
}

type CountAnomaliesResponseDTO struct {
	Current  *CountAnomaliesDTO `json:"current"`
	Previous *CountAnomaliesDTO `json:"previous"`
	Success  bool               `json:"success"`
	Error    string             `json:"error,omitempty"`
}
