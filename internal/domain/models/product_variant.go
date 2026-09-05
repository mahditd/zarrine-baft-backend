package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID uint `gorm:"primaryKey"`

	ProductID uint     `gorm:"uniqueIndex:idx_product_color_size"`
	Product   *Product `gorm:"foreignKey:ProductID"`

	ColorID uint   `gorm:"uniqueIndex:idx_product_color_size"`
	Color   *Color `gorm:"foreignKey:ColorID"`

	SizeID uint  `gorm:"uniqueIndex:idx_product_color_size"`
	Size   *Size `gorm:"foreignKey:SizeID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
