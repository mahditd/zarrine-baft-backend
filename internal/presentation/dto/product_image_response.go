package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductImageResponse struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	ImageURL  string `json:"image_url"`
}

func FromProductImage(
	image *models.ProductImage,
) ProductImageResponse {

	return ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		ImageURL:  image.ImageURL,
	}
}