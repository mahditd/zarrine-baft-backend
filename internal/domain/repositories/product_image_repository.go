package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductImageRepository interface {
	Create(image *models.ProductImage) error

	FindByID(id uint) (*models.ProductImage, error)

	FindByProductID(productID uint) ([]models.ProductImage, error)

	CountByProductID(productID uint) (int64, error)

	Update(image *models.ProductImage) error

	Reorder(productID uint, imageIDs []uint) error

	Delete(image *models.ProductImage) error
}
