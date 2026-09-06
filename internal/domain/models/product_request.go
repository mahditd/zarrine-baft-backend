package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductRequestStatus string

const (
	RequestNew          ProductRequestStatus = "new"
	RequestContacted    ProductRequestStatus = "contacted"
	RequestInDiscussion ProductRequestStatus = "in_discussion"
	RequestCompleted    ProductRequestStatus = "completed"
	RequestCancelled    ProductRequestStatus = "cancelled"
)

type ProductRequest struct {
	ID uint `gorm:"primaryKey"`

	RequestNumber string `gorm:"type:varchar(20);uniqueIndex;not null"`

	UserID uint  `gorm:"not null;index"`
	User   *User `gorm:"foreignKey:UserID"`

	CustomerName string `gorm:"not null"`

	Phone string `gorm:"not null;index"`

	CompanyName string `gorm:"not null"`

	CompanyPhone string `gorm:"not null"`

	Description string

	Status ProductRequestStatus `gorm:"type:varchar(20);not null;default:'new'"`

	Items []ProductRequestItem `gorm:"foreignKey:RequestID"`

	StatusHistory []ProductRequestStatusHistory `gorm:"foreignKey:RequestID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
