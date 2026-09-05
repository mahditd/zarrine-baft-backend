package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductImageRepository interface {

	Create(
		image *models.ProductImage,
	) error

	FindByProductID(
		productID uint,
	) ([]models.ProductImage, error)

	Delete(
		id uint,
	) error
}