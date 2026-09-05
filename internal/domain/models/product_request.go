package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductRequest struct {
	ID uint `gorm:"primaryKey"`

	CustomerName string `gorm:"not null"`
	Phone        string `gorm:"not null"`

	Description string

	Status string `gorm:"not null;default:'pending'"`

	Items []ProductRequestItem `gorm:"foreignKey:RequestID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
