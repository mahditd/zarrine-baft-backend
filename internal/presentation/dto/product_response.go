package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)


type ProductResponse struct {
	ID uint `json:"id"`

	NameFA string `json:"name_fa"`
	NameEN string `json:"name_en"`

	Category CategoryResponse `json:"category"`
	Material MaterialResponse `json:"material"`
}


func FromProduct(
	product *models.Product,
) ProductResponse {

	return ProductResponse{
		ID: product.ID,

		NameFA: product.NameFA,
		NameEN: product.NameEN,

		Category: FromCategory(&product.Category),
		Material: FromMaterial(&product.Material),
	}
}