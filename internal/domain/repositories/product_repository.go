package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductFilter struct {
	Page        int
	Limit       int
	Search      string
	CategoryIDs []uint
	MaterialIDs []uint
	ColorIDs    []uint
	SizeIDs     []uint
	IsActive    *bool
}

type ProductRepository interface {
	Create(product *models.Product) error

	FindAll() ([]models.Product, error)

	FindByID(id uint) (*models.Product, error)

	FindByProductCode(code string) (*models.Product, error)

	FindActiveProducts(filter ProductFilter) ([]models.Product, int64, error)

	FindAdminProducts(filter ProductFilter) ([]models.Product, int64, error)

	FindActiveByID(id uint) (*models.Product, error)

	Update(product *models.Product) error

	Reorder(productIDs []uint) error

	Delete(product *models.Product) error

	GetActiveCount() (int64, error)

	GetInactiveCount() (int64, error)
}
