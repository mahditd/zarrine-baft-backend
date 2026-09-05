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

func (r *ProductRepositoryImpl) Update(
	product *models.Product,
) error {

	return r.db.Save(product).Error
}

func (r *ProductRepositoryImpl) Delete(
	product *models.Product,
) error {

	return r.db.Delete(product).Error
}

func (r *ProductRepositoryImpl) FindActiveProducts(
	page int,
	limit int,
	categoryID uint,
	materialID uint,
	colorID uint,
) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	query := r.db.
		Model(&models.Product{}).
		Where("is_active = ?", true)

	if categoryID != 0 {

		query = query.Where(
			"category_id = ?",
			categoryID,
		)

	}

	if materialID != 0 {

		query = query.Where(
			"material_id = ?",
			materialID,
		)

	}

	if colorID != 0 {

		query = query.
			Joins(
				"JOIN product_variants ON product_variants.product_id = products.id",
			).
			Where(
				"product_variants.color_id = ?",
				colorID,
			).
			Group(
				"products.id",
			)

	}

	err := query.
		Count(&total).
		Error

	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = query.
		Preload("Category").
		Preload("Material").
		Preload("Images").
		Preload("Variants.Color").
		Preload("Variants.Size").
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error

	return products, total, err
}
func (r *ProductRepositoryImpl) FindActiveByID(
	id uint,
) (*models.Product, error) {

	var product models.Product

	err := r.db.
		Preload("Category").
		Preload("Material").
		Preload("Images").
		Preload("Variants.Color").
		Preload("Variants.Size").
		Where("is_active = ?", true).
		First(&product, id).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}
