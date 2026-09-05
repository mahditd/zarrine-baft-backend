package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID uint `gorm:"primaryKey"`

	ProductID uint `gorm:"uniqueIndex:idx_product_color_size"`
	Product   *Product

	ColorID uint `gorm:"uniqueIndex:idx_product_color_size"`
	Color   *Color

	SizeID uint `gorm:"uniqueIndex:idx_product_color_size"`
	Size   *Size

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
