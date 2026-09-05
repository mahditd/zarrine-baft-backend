package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type MaterialRepository interface {
	Create(material *models.Material) error

	FindByNameFA(name string) (*models.Material, error)

	FindByNameEN(name string) (*models.Material, error)

	FindByID(id uint) (*models.Material, error)

	FindAll() ([]models.Material, error)
}
