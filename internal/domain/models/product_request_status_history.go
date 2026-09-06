package models

import "time"

type ProductRequestStatusHistory struct {
	ID uint `gorm:"primaryKey"`

	RequestID uint `gorm:"not null;index"`

	FromStatus ProductRequestStatus `gorm:"type:varchar(20);not null"`

	ToStatus ProductRequestStatus `gorm:"type:varchar(20);not null"`

	Note string

	AdminID *uint `gorm:"index"`
	Admin   *User `gorm:"foreignKey:AdminID"`

	CreatedAt time.Time
}