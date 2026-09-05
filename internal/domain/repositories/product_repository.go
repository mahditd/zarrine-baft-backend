package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductRepository interface {
	Create(product *models.Product) error

	FindAll() ([]models.Product, error)

	FindByID(id uint) (*models.Product, error)

	FindActiveProducts(
		page int,
		limit int,
		categoryID uint,
		materialID uint,
		colorID uint,
	) ([]models.Product, int64, error)

	FindActiveByID(id uint) (*models.Product, error)

	Update(product *models.Product) error

	Delete(product *models.Product) error
}
