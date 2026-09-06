package repositories

import (
	"errors"
	"fmt"

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

func (r *ProductRequestRepositoryImpl) CreateStatusHistory(
	history *models.ProductRequestStatusHistory,
) error {

	return r.db.Create(history).Error
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
		Order("created_at DESC").
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
		Preload("User").
		Preload("Items.ProductVariant.Product").
		Preload("Items.ProductVariant.Color").
		Preload("Items.ProductVariant.Size").
		Preload("StatusHistory.Admin").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_request_status_histories.created_at ASC")
		})
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

	err = preloadProductRequest(query).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&requests).
		Error

	return requests, total, err
}

func (r *ProductRequestRepositoryImpl) FindByUserIDPaginated(
	userID uint,
	page int,
	limit int,
) ([]models.ProductRequest, int64, error) {

	var requests []models.ProductRequest
	var total int64

	query := r.db.Model(&models.ProductRequest{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = preloadProductRequest(query).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&requests).
		Error

	return requests, total, err
}

func (r *ProductRequestRepositoryImpl) FindByIDAndUserID(
	id uint,
	userID uint,
) (*models.ProductRequest, error) {

	var request models.ProductRequest

	err := preloadProductRequest(r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&request).
		Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *ProductRequestRepositoryImpl) GetLatestRequestNumber(
	year int,
) (string, error) {

	prefix := fmt.Sprintf("%04d-%%", year)
	var latest models.ProductRequest

	err := r.db.
		Unscoped().
		Where("request_number LIKE ?", prefix).
		Order("request_number DESC").
		First(&latest).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return latest.RequestNumber, nil
}

