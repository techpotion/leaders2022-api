package dto

import "gitlab.com/techpotion/leadershack2022/api/entity"

type GetRequestsFullRequestDTO struct {
	GetPointsRequestQueryDTO
	AnomalyCases []int `form:"-"`
}

type GetRequestsFullResponseDTO struct {
	Requests []*entity.RequestFull `json:"requests"`
	Count    int                   `json:"count"`
	Success  bool                  `json:"success"`
	Error    string                `json:"error,omitempty"`
}
