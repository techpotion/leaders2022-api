package dto

import (
	"time"

	"gitlab.com/techpotion/leadershack2022/api/entity"
)

type CountAnomaliesGrouppedRequestDTO struct {
	Region       string    `form:"region"          validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type CountAnomaliesGrouppedResponseDTO struct {
	GroupsCounts []*entity.CountAnomaliesGroupped `json:"groups_counts"`
	Success      bool                             `json:"success"`
	Error        string                           `json:"error,omitempty"`
}

type CountAnomaliesByOwnerCompaniesRequestDTO struct {
	Region       string    `form:"region"          validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type CountAnomaliesByOwnerCompaniesResponseDTO struct {
	OwnerCompanies []*entity.AnomaliesByOwnerCompany `json:"owner_companies"`
	Count          int                               `json:"count"`
	Success        bool                              `json:"success"`
	Error          string                            `json:"error,omitempty"`
}

type CountAnomaliesByServingCompaniesRequestDTO struct {
	Region       string    `form:"region"          validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type CountAnomaliesByServingCompaniesResponseDTO struct {
	ServingCompanies []*entity.AnomaliesByServingCompany `json:"serving_companies"`
	Count            int                                 `json:"count"`
	Success          bool                                `json:"success"`
	Error            string                              `json:"error,omitempty"`
}

type CountAnomaliesByDeffectCategoriesRequestDTO struct {
	Region       string    `form:"region"          validate:"required"`
	DateTimeFrom time.Time `form:"datetime_from"   validate:"required"`
	DateTimeTo   time.Time `form:"datetime_to"     validate:"required"`
}

type CountAnomaliesByDeffectCategoriesResponseDTO struct {
	DeffectCategories []*entity.AnomaliesByDeffectCategory `json:"deffect_categories"`
	Count             int                                  `json:"count"`
	Success           bool                                 `json:"success"`
	Error             string                               `json:"error,omitempty"`
}
