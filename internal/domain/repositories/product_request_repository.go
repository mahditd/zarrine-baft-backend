package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductRequestRepository interface {
	Create(
		request *models.ProductRequest,
	) error

	FindAll() (
		[]models.ProductRequest,
		error,
	)

	FindPaginated(
		page int,
		limit int,
		status string,
	) (
		[]models.ProductRequest,
		int64,
		error,
	)

	FindByID(
		id uint,
	) (*models.ProductRequest, error)

	Update(
		request *models.ProductRequest,
	) error
}
