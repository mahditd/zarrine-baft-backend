package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type CategoryResponse struct {
	ID     uint   `json:"id"`
	NameFA string `json:"name_fa"`
	NameEN string `json:"name_en"`
}

func FromCategory(category *models.Category) CategoryResponse {

	return CategoryResponse{
		ID:     category.ID,
		NameFA: category.NameFA,
		NameEN: category.NameEN,
	}
}
