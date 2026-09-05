package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ColorRepository interface {
	Create(color *models.Color) error

	FindByNameFA(name string) (*models.Color, error)

	FindByNameEN(name string) (*models.Color, error)

	FindByID(id uint) (*models.Color, error)

	FindAll() ([]models.Color, error)
}
