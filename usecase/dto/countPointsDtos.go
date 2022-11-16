package dto

import "time"

type CountPointsRequestDTO struct {
	Region          string    `form:"region"          validate:"required"` // район
	DateTimeFrom    time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo      time.Time `form:"datetime_to"     validate:"required"`
	XMin            *float64  `form:"x_min"           validate:"gte=-180,lte=180"`
	Ymin            *float64  `form:"y_min"           validate:"gte=-90,lte=90"`
	XMax            *float64  `form:"x_max"           validate:"gte=-180,lte=180"`
	YMax            *float64  `form:"y_max"           validate:"gte=-90,lte=90"`
	ServingCompany  []string  `form:"serving_company"`   // компания исполнитель
	Efficiency      []string  `form:"efficiency"`        // результативность
	GradeForService []string  `form:"grade_for_service"` // оценка качества
	UrgencyCategory []string  `form:"urgency_category"`  // срочность
	WorkType        []string  `form:"work_type"`         // виды работ
	DeffectCategory []string  `form:"deffect_category"`  // категория дефекта
	OwnerCompany    []string  `form:"owner_company"`     // упр. организация
}

type CountPointsResponseDTO struct {
	Count   int    `json:"count"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
