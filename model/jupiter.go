package model

import (
	"time"

	"gorm.io/gorm"
)

type JupiterData struct {
	ID            int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	UID           int       `json:"uid,omitempty"`
	Account       string    `json:"account,omitempty"`
	Password      string    `json:"password,omitempty"`
	GPA           float32   `json:"gpa,omitempty"`
	DataFetchedAt time.Time `json:"updated_at,omitempty"`
}

type Course struct {
	ID           int    `json:"id,omitempty" gorm:"primaryKey,autoIncrement"`
	UID          int    `json:"uid,omitempty"`
	Title        string `json:"title,omitempty"`
	PercentGrade string `json:"percent_grade,omitempty"`
	LetterGrade  string `json:"letter_grade,omitempty"`
}

type Assignment struct {
	ID     int       `json:"id,omitempty" gorm:"primaryKey,autoIncrement"`
	UID    int       `json:"uid,omitempty"`
	Status int       `json:"status,omitempty"` // 常规 [0], 完成 [1], 检索相关数据中 [-1]
	From   string    `json:"from,omitempty"`
	Due    time.Time `json:"due,omitempty"`
	Title  string    `json:"title,omitempty"`
	Desc   string    `json:"desc,omitempty"`
	Score  string    `json:"score,omitempty"`
}

func (o *Course) BeforeUpdate(db *gorm.DB) error {
	if db.Statement.Changed("PercentGrade") || db.Statement.Changed("LetterGrade") {
		println(o.Title + " has a change in grade")
	}
	return nil
}
