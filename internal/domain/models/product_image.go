package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	ID uint `gorm:"primaryKey"`

	ProductID uint `gorm:"not null"`
	Product   Product

	ImageURL string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
