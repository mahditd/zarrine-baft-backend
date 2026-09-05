package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(
	db *gorm.DB,
) domainRepositories.ProductRepository {

	return &ProductRepositoryImpl{
		db: db,
	}
}

func (r *ProductRepositoryImpl) Create(
	product *models.Product,
) error {

	return r.db.Create(product).Error
}

func (r *ProductRepositoryImpl) FindByNameFA(
	name string,
) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Where("name_fa = ?", name).
		First(&product).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) FindByNameEN(
	name string,
) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Where("name_en = ?", name).
		First(&product).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) FindAll() ([]models.Product, error) {

	var products []models.Product

	err := r.db.
		Preload("Category").
		Preload("Material").
		Preload("Images").
		Preload("Variants.Color").
		Preload("Variants.Size").
		Find(&products).
		Error

	return products, err
}

func (r *ProductRepositoryImpl) FindByID(
	id uint,
) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Preload("Category").
		Preload("Material").
		Preload("Images").
		Preload("Variants.Color").
		Preload("Variants.Size").
		First(&product, id).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}
