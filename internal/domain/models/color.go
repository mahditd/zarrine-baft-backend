package models

import (
	"time"

	"gorm.io/gorm"
)

type Color struct {
	ID uint `gorm:"primaryKey"`

	NameFA string `gorm:"not null"`
	NameEN string `gorm:"not null"`

	HexCode string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}