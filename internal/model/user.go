package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Status    int       `json:"status,omitempty"`
	Email     string    `json:"email,omitempty" gorm:"unique"`
	Username  string    `json:"username,omitempty" gorm:"unique"`
	Password  string    `json:"password,omitempty"`
	ActiveAt  time.Time `json:"active_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
