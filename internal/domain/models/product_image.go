package models

import (
	"time"
)

type ProductImage struct {
	ID uint `gorm:"primaryKey"`

	ProductID uint     `gorm:"not null;index"`
	Product   *Product `gorm:"foreignKey:ProductID"`

	ImageURL string `gorm:"not null"`
	FilePath string `gorm:"not null"`

	DisplayOrder int  `gorm:"not null;default:1"`
	IsCover      bool `gorm:"not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
