package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductRequestItem struct {
	ID uint `gorm:"primaryKey"`

	RequestID uint
	Request   ProductRequest `gorm:"foreignKey:RequestID"`

	ProductVariantID uint
	ProductVariant   ProductVariant

	Quantity int `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
