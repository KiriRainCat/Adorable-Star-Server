package model

import (
	"time"
)

type Message struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey,autoIncrement"`
	UID       int       `json:"uid,omitempty"`
	Type      int       `json:"type,omitempty"` // 系统通知 [-1], 新作业 [0], 作业信息变动 [1], 科目成绩变动 [2]
	From      int       `json:"from,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
