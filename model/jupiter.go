package model

import (
	"time"

	"gorm.io/gorm"
)

type JupiterData struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	UID       int       `json:"uid,omitempty"`
	Account   string    `json:"account,omitempty"`
	Password  string    `json:"password,omitempty"`
	GPA       string    `json:"gpa,omitempty"`
	FetchedAt time.Time `json:"fetched_at,omitempty"`
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

// Copy all fields from [other] to this course for EMPTY fields
func (o *Course) CopyFromOther(other *Course) {
	if o.ID == 0 {
		o.ID = other.ID
	}
	if o.UID == 0 {
		o.UID = other.UID
	}
	if o.Title == "" {
		o.Title = other.Title
	}
	if o.PercentGrade == "" {
		o.PercentGrade = other.PercentGrade
	}
	if o.LetterGrade == "" {
		o.LetterGrade = other.PercentGrade
	}
}

// Copy all fields from [other] to this assignment for EMPTY fields
func (o *Assignment) CopyFromOther(other *Assignment) {
	if o.ID == 0 {
		o.ID = other.ID
	}
	if o.UID == 0 {
		o.UID = other.UID
	}
	if o.Status == 0 {
		o.Status = other.Status
	}
	if o.From == "" {
		o.From = other.From
	}
	if (o.Due == time.Time{}) {
		o.Due = other.Due
	}
	if o.Title == "" {
		o.Title = other.Title
	}
	if o.Desc == "" {
		o.Desc = other.Desc
	}
	if o.Score == "" {
		o.Score = other.Score
	}
}

func (o *Course) BeforeUpdate(db *gorm.DB) (err error) {
	if db.Statement.Changed("PercentGrade") || db.Statement.Changed("LetterGrade") {
		db.Create(&Message{
			UID:  o.UID,
			Type: 1,
			From: o.Title,
			Msg:  "新成绩: [" + o.PercentGrade + " " + o.LetterGrade + "]",
		})
	}
	return nil
}
