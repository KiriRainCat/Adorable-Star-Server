package model

import (
	"time"
)

type Message struct {
	// 新作业 [0], 新成绩 [1], 科目成绩变动 [2]
	ID        int       `json:"id,omitempty" gorm:"primaryKey,autoIncrement"`
	UID       int       `json:"uid,omitempty"`
	Type      int       `json:"type,omitempty"`
	From      string    `json:"from,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
