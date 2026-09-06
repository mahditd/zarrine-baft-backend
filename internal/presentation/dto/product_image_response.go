package dto

import "github.com/mahditd/zarrine-baft-backend/internal/domain/models"

type ProductImageResponse struct {
	ID           uint   `json:"id"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
	IsCover      bool   `json:"is_cover"`
}

func FromProductImage(
	image *models.ProductImage,
) ProductImageResponse {

	return ProductImageResponse{
		ID:           image.ID,
		ImageURL:     image.ImageURL,
		DisplayOrder: image.DisplayOrder,
		IsCover:      image.IsCover,
	}
}
