package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductVariantRepository interface {
	Create(
		variant *models.ProductVariant,
	) error

	FindByID(
		id uint,
	) (*models.ProductVariant, error)

	FindByProductID(
		productID uint,
	) ([]models.ProductVariant, error)

	FindByProductColorAndSize(
		productID uint,
		colorID uint,
		sizeID uint,
	) (*models.ProductVariant, error)

	FindAll() (
		[]models.ProductVariant,
		error,
	)
}
