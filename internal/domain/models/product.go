package models

import (
	"time"
)

type Product struct {
	ID uint `gorm:"primaryKey"`

	ProductCode string `gorm:"uniqueIndex;not null"`

	NameFA string `gorm:"not null"`

	NameEN string `gorm:"not null"`

	IsActive bool `gorm:"not null;default:true"`

	DisplayOrder int `gorm:"not null;default:0"`

	CategoryID uint `gorm:"not null"`

	Category Category

	MaterialID uint `gorm:"not null"`

	Material Material

	Images []ProductImage `gorm:"foreignKey:ProductID"`

	Variants []ProductVariant `gorm:"foreignKey:ProductID"`

	CreatedAt time.Time

	UpdatedAt time.Time
}
