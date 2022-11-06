package dto

type GetRegionsResponseDTO struct {
	Regions []string `json:"regions"`
	Count   int      `json:"count"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}
