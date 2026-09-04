package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductVariantResponse struct {
	ID uint `json:"id"`

	ProductID uint `json:"product_id"`

	Color ColorResponse `json:"color"`
}

func FromProductVariant(
	variant *models.ProductVariant,
) ProductVariantResponse {

	return ProductVariantResponse{
		ID: variant.ID,

		ProductID: variant.ProductID,

		Color: FromColor(variant.Color),
	}
}
