package entity

type RequestFull struct {
	RootID string `json:"root_id"`
	*Request
	*HCSPoint
	*RequestAnomaly
}
