package dto

type GetRegionsResponseDTO struct {
	Regions []string `json:"regions"`
	Count   int      `json:"count"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}

type GetServingCompaniesResponseDTO struct {
	ServingCompanies []string `json:"serving_companies"`
	Count            int      `json:"count"`
	Success          bool     `json:"success"`
	Error            string   `json:"error,omitempty"`
}

type GetOwnerCompaniesResponseDTO struct {
	OwnerCompanies []string `json:"owner_companies"`
	Count          int      `json:"count"`
	Success        bool     `json:"success"`
	Error          string   `json:"error,omitempty"`
}

type GetDeffectCategoriesResponseDTO struct {
	DeffectCategories []string `json:"deffect_categories"`
	Count             int      `json:"count"`
	Success           bool     `json:"success"`
	Error             string   `json:"error,omitempty"`
}

type GetWorkTypesResponseDTO struct {
	WorkTypes []string `json:"work_types"`
	Count     int      `json:"count"`
	Success   bool     `json:"success"`
	Error     string   `json:"error,omitempty"`
}

type GetDispatchersResponseDTO struct {
	Dispatchers []string `json:"dispatcher"`
	Count       int      `json:"count"`
	Success     bool     `json:"success"`
	Error       string   `json:"error,omitempty"`
}
