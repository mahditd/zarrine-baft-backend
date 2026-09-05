package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type ProductRequestRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRequestRepository(
	db *gorm.DB,
) domainRepositories.ProductRequestRepository {

	return &ProductRequestRepositoryImpl{
		db: db,
	}
}

func (r *ProductRequestRepositoryImpl) Update(
	request *models.ProductRequest,
) error {

	return r.db.Save(request).Error
}

func (r *ProductRequestRepositoryImpl) Create(
	request *models.ProductRequest,
) error {

	return r.db.Create(request).Error
}

func (r *ProductRequestRepositoryImpl) FindAll() (
	[]models.ProductRequest,
	error,
) {

	var requests []models.ProductRequest

	err := preloadProductRequest(r.db).
		Find(&requests).
		Error

	return requests, err
}

func (r *ProductRequestRepositoryImpl) FindByID(
	id uint,
) (*models.ProductRequest, error) {

	var request models.ProductRequest

	err := preloadProductRequest(r.db).
		First(&request, id).
		Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func preloadProductRequest(db *gorm.DB) *gorm.DB {

	return db.
		Preload("Items.ProductVariant.Product").
		Preload("Items.ProductVariant.Color").
		Preload("Items.ProductVariant.Size")
}

func (r *ProductRequestRepositoryImpl) FindPaginated(
	page int,
	limit int,
	status string,
) ([]models.ProductRequest, int64, error) {

	var requests []models.ProductRequest
	var total int64

	query := r.db.Model(&models.ProductRequest{})

	if status != "" {

		query = query.Where(
			"status = ?",
			status,
		)

	}

	err := query.Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = query.
		Preload("Items").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&requests).
		Error

	return requests, total, err
}
