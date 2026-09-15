package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductRequestItem struct {
	ID uint `gorm:"primaryKey"`

	RequestID uint
	Request   *ProductRequest `gorm:"foreignKey:RequestID"`

	ProductVariantID uint
	ProductVariant   *ProductVariant `gorm:"foreignKey:ProductVariantID"`

	Quantity int `gorm:"not null"`
	PriceSnapshot int64 `gorm:"not null"`

	// Snapshot fields (SRS 12.2 & 14): preserve request data even if product/variant is edited or deleted later
	ProductCode   string `gorm:"type:varchar(10)"`
	ProductNameFA string `gorm:"type:varchar(255)"`
	ProductNameEN string `gorm:"type:varchar(255)"`
	ColorNameFA   string `gorm:"type:varchar(100)"`
	ColorNameEN   string `gorm:"type:varchar(100)"`
	SizeName      string `gorm:"type:varchar(20)"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
