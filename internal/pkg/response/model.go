package response

import (
	"adorable-star/internal/model"
	"time"
)

type Course struct {
	model.Course
	Assignments []*model.Assignment `json:"assignments,omitempty"`
}

type Data struct {
	FetchedAt time.Time `json:"fetched_at,omitempty"`
	GPA       string    `json:"gpa,omitempty"`
	Data      any       `json:"data,omitempty"`
}
