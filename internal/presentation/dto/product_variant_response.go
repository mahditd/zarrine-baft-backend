package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type SizeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductVariantResponse struct {
	ID        uint          `json:"id"`
	ProductID uint          `json:"product_id"`
	Color     ColorResponse `json:"color"`
	Size      SizeResponse  `json:"size"`
}

func FromProductVariant(
	variant *models.ProductVariant,
) ProductVariantResponse {

	return ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Color:     FromColor(variant.Color),
		Size: SizeResponse{
			ID:   variant.Size.ID,
			Name: variant.Size.Name,
		},
	}
}
