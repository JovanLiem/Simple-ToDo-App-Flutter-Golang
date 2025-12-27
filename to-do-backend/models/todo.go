package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title" binding:"required"`
	IsCompleted bool           `json:"is_completed" gorm:"default:false"`
	Date        time.Time      `json:"date"`
	Deadline    *time.Time     `json:"deadline"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
