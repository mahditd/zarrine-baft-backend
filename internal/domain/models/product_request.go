package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductRequestStatus string

const (
	RequestPending   ProductRequestStatus = "pending"
	RequestReviewing ProductRequestStatus = "reviewing"
	RequestApproved  ProductRequestStatus = "approved"
	RequestRejected  ProductRequestStatus = "rejected"
	RequestCompleted ProductRequestStatus = "completed"
)

type ProductRequest struct {
	ID uint `gorm:"primaryKey"`

	CustomerName string `gorm:"not null"`

	Phone string `gorm:"not null;index"`

	CompanyName string `gorm:"not null"`

	CompanyPhone string `gorm:"not null"`

	Description string

	Status ProductRequestStatus `gorm:"type:varchar(20);not null;default:'pending'"`

	Items []ProductRequestItem `gorm:"foreignKey:RequestID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}