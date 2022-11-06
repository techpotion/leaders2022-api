package entity

type RequestAnomaly struct {
	RootID         string   `json:"root_id"`
	IsAnomaly      bool     `json:"is_anomaly"`
	AnomalyCases   []int    `json:"anomaly_cases"`
	IsCustom       bool     `json:"is_custom"`
	NetProbability *float32 `json:"net_probability"`
}
