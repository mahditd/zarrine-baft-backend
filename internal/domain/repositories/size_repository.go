package repositories

import "github.com/mahditd/zarrine-baft-backend/internal/domain/models"

type SizeRepository interface {

	FindByID(
		id uint,
	) (*models.Size, error)

}