package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey"`

	NameFA string `gorm:"not null"`
	NameEN string `gorm:"not null"`

	IsActive bool `gorm:"not null;default:true"`

	CategoryID uint `gorm:"not null"`
	Category   Category

	MaterialID uint `gorm:"not null"`
	Material   Material

	Images []ProductImage `gorm:"foreignKey:ProductID"`

	Variants []ProductVariant `gorm:"foreignKey:ProductID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
