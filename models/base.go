package models

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        int64          `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` //add soft delete in gorm
}
