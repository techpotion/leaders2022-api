package dto

import "time"

type CountRequestsFullRequestDTO struct {
	Region          string    `form:"region"          validate:"required"`
	DateTimeFrom    time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo      time.Time `form:"datetime_to"     validate:"required"`
	XMin            *float64  `form:"x_min"           validate:"gte=-180,lte=180"`
	Ymin            *float64  `form:"y_min"           validate:"gte=-90,lte=90"`
	XMax            *float64  `form:"x_max"           validate:"gte=-180,lte=180"`
	YMax            *float64  `form:"y_max"           validate:"gte=-90,lte=90"`
	UrgencyCategory *string   `form:"urgency_category"`
	AnomalyCases    []int     `form:"-"`
}

type CountRequestsFullResponseDTO struct {
	Count   int    `json:"count"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
