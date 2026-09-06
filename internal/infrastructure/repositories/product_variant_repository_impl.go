package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type ProductVariantRepositoryImpl struct {
	db *gorm.DB
}

func NewProductVariantRepository(
	db *gorm.DB,
) domainRepositories.ProductVariantRepository {

	return &ProductVariantRepositoryImpl{
		db: db,
	}
}

func (r *ProductVariantRepositoryImpl) Create(
	variant *models.ProductVariant,
) error {

	return r.db.Create(variant).Error
}

func (r *ProductVariantRepositoryImpl) FindByProductID(
	productID uint,
) ([]models.ProductVariant, error) {

	var variants []models.ProductVariant

	err := r.db.
		Preload("Color").
		Preload("Size").
		Where("product_id = ?", productID).
		Find(&variants).
		Error

	return variants, err
}

func (r *ProductVariantRepositoryImpl) FindAll() ([]models.ProductVariant, error) {

	var variants []models.ProductVariant

	err := r.db.
		Preload("Product").
		Preload("Color").
		Preload("Size").
		Find(&variants).
		Error

	return variants, err
}

func (r *ProductVariantRepositoryImpl) FindByID(
	id uint,
) (*models.ProductVariant, error) {

	var variant models.ProductVariant

	err := r.db.
		Preload("Product").
		Preload("Color").
		Preload("Size").
		First(&variant, id).
		Error

	if err != nil {
		return nil, err
	}

	return &variant, nil
}
func (r *ProductVariantRepositoryImpl) FindByProductColorAndSize(
	productID uint,
	colorID uint,
	sizeID uint,
) (*models.ProductVariant, error) {

	var variant models.ProductVariant

	err := r.db.
		Where(
			"product_id = ? AND color_id = ? AND size_id = ?",
			productID,
			colorID,
			sizeID,
		).
		First(&variant).
		Error

	if err != nil {
		return nil, err
	}

	return &variant, nil
}

func (r *ProductVariantRepositoryImpl) Update(
	variant *models.ProductVariant,
) error {

	return r.db.Save(variant).Error
}

func (r *ProductVariantRepositoryImpl) Delete(
	variant *models.ProductVariant,
) error {

	return r.db.Delete(variant).Error
}
