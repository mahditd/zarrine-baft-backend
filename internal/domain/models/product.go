package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey"`

	NameFA string `gorm:"not null"`
	NameEN string `gorm:"not null"`

	CategoryID uint
	Category Category

	MaterialID uint
	Material Material

	CreatedAt time.Time
	UpdatedAt time.Time

	Images []ProductImage

	DeletedAt gorm.DeletedAt `gorm:"index"`
}