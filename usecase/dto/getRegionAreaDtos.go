package dto

type GetRegionAreaResponseDTO struct {
	AreaPloygonGeoJSON string `json:"area_polygon_geojson"`
	Success            bool   `json:"success"`
	Error              string `json:"error,omitempty"`
}
