package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductRepository interface {

	Create(product *models.Product) error

	FindAll() ([]models.Product, error)

	FindByID(id uint) (*models.Product, error)
}