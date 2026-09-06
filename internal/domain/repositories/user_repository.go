package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type UserRepository interface {
	Create(user *models.User) error

	FindByPhone(phone string) (*models.User, error)

	FindByEmail(email string) (*models.User, error)

	FindByID(id uint) (*models.User, error)

	Update(user *models.User) error
}
