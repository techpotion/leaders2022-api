package dto

import "gitlab.com/techpotion/leadershack2022/api/entity"

type GetRequestsByIdsRequestDto struct {
	RootIds string `form:"root_ids" validate:"required"`
}

type GetRequestsByIdsResponseDto struct {
	Requests []*entity.Request `json:"requests"`
	Count    int               `json:"count"`
	Success  bool              `json:"success"`
	Error    string            `json:"error,omitempty"`
}
