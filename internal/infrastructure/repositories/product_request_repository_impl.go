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
