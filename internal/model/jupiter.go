package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type JupiterData struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	UID       int       `json:"uid,omitempty" gorm:"unique"`
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
		o.LetterGrade = other.LetterGrade
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
	if o.Due.Year() == (time.Time{}.Year()) {
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

func (o *Course) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("PercentGrade", "LetterGrade") {
		percentGrade := o.PercentGrade
		letterGrade := o.LetterGrade

		go func() {
			// Wait for course to be updated
			time.Sleep(time.Second * 6)

			// Insert new message to database
			tx.Create(&Message{
				UID:    o.UID,
				Type:   2,
				From:   o.ID,
				Course: o.Title,
				Msg:    percentGrade + " " + letterGrade + "|" + o.PercentGrade + " " + o.LetterGrade,
			})
		}()
	}
	return nil
}

func (o *Assignment) BeforeUpdate(tx *gorm.DB) error {
	// If due date update
	if tx.Statement.Changed("Due") {
		due := o.Due
		go func() {
			// Wait for course to be updated
			time.Sleep(time.Second * 6)

			// Insert new message to database
			date := "Future"
			if o.Due.Year() != 1 {
				date = strconv.Itoa(int(due.Month())) + "/" + strconv.Itoa(due.Day())
			}
			tx.Create(&Message{
				UID:    o.UID,
				Type:   1,
				From:   o.ID,
				Course: o.From,
				Msg:    "Due|" + date + "|" + strconv.Itoa(int(o.Due.Month())) + "/" + strconv.Itoa(o.Due.Day()),
			})
		}()
	}

	// If score update
	if tx.Statement.Changed("Score") {
		score := o.Score
		go func() {
			// Wait for course to be updated
			time.Sleep(time.Second * 6)

			// Insert new message to database
			tx.Create(&Message{
				UID:    o.UID,
				Type:   1,
				From:   o.ID,
				Course: o.From,
				Msg:    "Score|" + score + "|" + o.Score,
			})
		}()
	}

	// If description update
	if tx.Statement.Changed("Desc") {
		desc := o.Desc
		go func() {
			// Wait for course to be updated
			time.Sleep(time.Second * 6)

			// Insert new message to database
			tx.Create(&Message{
				UID:    o.UID,
				Type:   1,
				From:   o.ID,
				Course: o.From,
				Msg:    "Desc|" + desc + "|" + o.Desc,
			})
		}()
	}
	return nil
}

func (o *Assignment) AfterCreate(tx *gorm.DB) error {
	// Check whether user is new user
	var data *JupiterData
	if tx.Where("uid = ?", o.UID).First(&data).Error != nil {
		return nil
	}
	if (data.FetchedAt == time.Time{}) {
		return nil
	}

	// Insert new message to database
	due := "Future"
	if o.Due.Year() != 1 {
		due = strconv.Itoa(int(o.Due.Month())) + "/" + strconv.Itoa(o.Due.Day())
	}
	tx.Create(&Message{
		UID:    o.UID,
		Type:   0,
		From:   o.ID,
		Course: o.From,
		Msg:    due + "|" + o.Title,
	})
	return nil
}
