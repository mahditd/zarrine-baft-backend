package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	FullName string `gorm:"not null"`

	Phone string `gorm:"uniqueIndex;not null"`

	Email *string `gorm:"uniqueIndex;type:varchar(255)"`

	PasswordHash string `gorm:"not null"`

	Role UserRole `gorm:"type:varchar(20);not null;default:'customer'"`

	CompanyName  *string `gorm:"type:varchar(100)"`
	CompanyPhone *string `gorm:"type:varchar(20)"`

	Country *string
	Address *string

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
