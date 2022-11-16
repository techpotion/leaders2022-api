package dto

import "time"

type GetEfficiencyPlotRequestDTO struct {
	Dispatcher   string    `form:"dispatcher"      validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type GetEfficiencyPlotResponseDTO struct {
	Filename *string `json:"filename"`
	Success  bool    `json:"success"`
	Error    string  `json:"error,omitempty"`
}
