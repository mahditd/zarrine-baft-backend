package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductResponse struct {
	ID uint `json:"id"`

	NameFA string `json:"name_fa"`
	NameEN string `json:"name_en"`

	IsActive bool `json:"is_active"`

	Category CategoryResponse `json:"category"`
	Material MaterialResponse `json:"material"`

	Images   []ProductImageResponse   `json:"images"`
	Variants []ProductVariantResponse `json:"variants"`
}

func FromProduct(
	product *models.Product,
) ProductResponse {

	images := make([]ProductImageResponse, 0)

	for _, image := range product.Images {
		images = append(
			images,
			FromProductImage(&image),
		)
	}

	variants := make([]ProductVariantResponse, 0)

	for _, variant := range product.Variants {
		variants = append(
			variants,
			FromProductVariant(&variant),
		)
	}

	return ProductResponse{
		ID: product.ID,

		NameFA: product.NameFA,
		NameEN: product.NameEN,

		IsActive: product.IsActive,

		Category: FromCategory(&product.Category),
		Material: FromMaterial(&product.Material),

		Images:   images,
		Variants: variants,
	}
}
