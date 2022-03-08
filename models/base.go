package models

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uint64         `json:"id"`
	CreatedAt time.Time      `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` //add soft delete in gorm
}
