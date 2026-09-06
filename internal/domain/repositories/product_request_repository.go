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
		search string,
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

	CreateStatusHistory(
		history *models.ProductRequestStatusHistory,
	) error

	FindByUserIDPaginated(
		userID uint,
		page int,
		limit int,
	) (
		[]models.ProductRequest,
		int64,
		error,
	)

	FindByIDAndUserID(
		id uint,
		userID uint,
	) (*models.ProductRequest, error)

	GetLatestRequestNumber(
		year int,
	) (string, error)

	GetNewRequestsCount() (int64, error)

	FindLatest(limit int) ([]models.ProductRequest, error)
}
