package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID uint `gorm:"primaryKey"`

	NameFA string `gorm:"not null"`
	NameEN string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}