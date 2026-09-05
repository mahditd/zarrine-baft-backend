package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type ProductImageRepositoryImpl struct {
	db *gorm.DB
}

func NewProductImageRepository(
	db *gorm.DB,
) domainRepositories.ProductImageRepository {

	return &ProductImageRepositoryImpl{
		db: db,
	}
}


func (r *ProductImageRepositoryImpl) Create(
	image *models.ProductImage,
) error {

	return r.db.Create(image).Error
}


func (r *ProductImageRepositoryImpl) FindByProductID(
	productID uint,
) ([]models.ProductImage, error) {

	var images []models.ProductImage

	err := r.db.
		Where("product_id = ?", productID).
		Find(&images).
		Error

	return images, err
}


func (r *ProductImageRepositoryImpl) Delete(
	id uint,
) error {

	return r.db.Delete(
		&models.ProductImage{},
		id,
	).Error
}