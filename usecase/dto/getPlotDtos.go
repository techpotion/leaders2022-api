package dto

import "time"

type GetPlotRequestDTO struct {
	RootID            string     `json:"root_id"`
	DispatchersNumber string     `json:"dispetchers_number"`
	Efficiency        *string    `json:"efficiency"`
	ClosureDate       *time.Time `json:"closure_date"`
	DateOfReview      *time.Time `json:"date_of_review"`
	GradeForService   *string    `json:"grade_for_service"`
}
