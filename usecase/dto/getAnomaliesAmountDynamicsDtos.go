package dto

import (
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
)

type GetAnomaliesAmountDynamicsRequestDTO struct {
	Region       string    `form:"region"          validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type GetAnomaliesAmountDynamicsResponseDTO struct {
	Dynamics []*entity.AnomaliesDynamics `json:"dynamics"`
	Count    int                         `json:"count"`
	Success  bool                        `json:"success"`
	Error    string                      `json:"error,omitempty"`
}
