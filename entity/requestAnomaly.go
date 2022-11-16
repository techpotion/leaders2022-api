package entity

import "time"

type RequestAnomaly struct {
	RootID         string   `json:"root_id"`
	IsAnomaly      bool     `json:"is_anomaly"`
	AnomalyCases   []int    `json:"anomaly_cases"`
	IsCustom       bool     `json:"is_custom"`
	NetProbability *float32 `json:"net_probability"`
}

type AnomaliesDynamics struct {
	Day    time.Time `json:"day"`
	Amount int       `json:"amount"`
}

type AnomaliesByOwnerCompany struct {
	OwnerCompany string  `json:"owner_company"`
	Count        int     `json:"count"`
	Percent      float32 `json:"percent"`
}

type AnomaliesByServingCompany struct {
	ServingCompany string  `json:"serving_company"`
	Count          int     `json:"count"`
	Percent        float32 `json:"percent"`
}

type AnomaliesByDeffectCategory struct {
	DeffectCategory string  `json:"deffect_category"`
	Count           int     `json:"count"`
	Percent         float32 `json:"percent"`
}

type CountAnomaliesGroupped struct {
	Type    int     `json:"type"`
	Count   int     `json:"count"`
	Percent float32 `json:"percent"`
}
