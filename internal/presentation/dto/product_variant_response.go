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
	Size      *SizeResponse `json:"size"`
}

func FromProductVariant(
	variant *models.ProductVariant,
) ProductVariantResponse {

	response := ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
	}

	if variant.Color != nil {
		response.Color = FromColor(variant.Color)
	}

	if variant.Size != nil {
		size := FromSize(variant.Size)
		response.Size = &size
	}

	return response
}

func FromSize(
	size *models.Size,
) SizeResponse {

	return SizeResponse{
		ID:   size.ID,
		Name: size.Name,
	}
}
