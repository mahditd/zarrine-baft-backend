package dto

import (
	"time"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type UserResponse struct {
	ID           uint      `json:"id"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	Email        *string   `json:"email"`
	Role         string    `json:"role"`
	CompanyName  *string   `json:"company_name"`
	CompanyPhone *string   `json:"company_phone"`
	Country      string    `json:"country"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
}

func FromUser(user *models.User) UserResponse {

	response := UserResponse{
		ID:           user.ID,
		FullName:     user.FullName,
		Phone:        user.Phone,
		Email:        user.Email,
		Role:         string(user.Role),
		CompanyName:  user.CompanyName,
		CompanyPhone: user.CompanyPhone,
		Country:      "",
		Address:      "",
		CreatedAt:    user.CreatedAt,
	}

	if user.Country != nil {
		response.Country = *user.Country
	}

	if user.Address != nil {
		response.Address = *user.Address
	}

	return response
}
