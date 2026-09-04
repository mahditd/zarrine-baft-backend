package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type CategoryRepository interface {

	Create(category *models.Category) error

	FindByNameFA(name string) (*models.Category, error)

	FindByNameEN(name string) (*models.Category, error)

	FindByID(id uint) (*models.Category, error)

	FindAll() ([]models.Category, error)
}