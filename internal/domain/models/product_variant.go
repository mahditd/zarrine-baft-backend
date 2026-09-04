package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID uint `gorm:"primaryKey"`

	ProductID uint `gorm:"uniqueIndex:idx_product_color"`
	Product   *Product

	ColorID uint `gorm:"uniqueIndex:idx_product_color"`
	Color   *Color

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
