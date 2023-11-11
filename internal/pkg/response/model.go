package response

import "adorable-star/internal/model"

type Course struct {
	model.Course
	Assignments []*model.Assignment `json:"assignments,omitempty"`
}
