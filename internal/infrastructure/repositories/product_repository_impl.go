package repositories

import (
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/utils"

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

	return r.db.Transaction(func(tx *gorm.DB) error {
		// New products appear at the top (display_order = 1)
		if err := tx.Model(&models.Product{}).
			Where("display_order >= 1").
			UpdateColumn("display_order", gorm.Expr("display_order + 1")).Error; err != nil {
			return err
		}

		product.DisplayOrder = 1
		return tx.Create(product).Error
	})
}

func (r *ProductRepositoryImpl) FindByNameFA(
	name string,
) (*models.Product, error) {

	var product models.Product

	err := preloadProduct(r.db).
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

	err := preloadProduct(r.db).
		Where("LOWER(name_en) = LOWER(?)", name).
		First(&product).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) FindByProductCode(
	code string,
) (*models.Product, error) {

	var product models.Product

	err := preloadProduct(r.db).
		Where("product_code = ?", code).
		First(&product).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func preloadProduct(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Category").
		Preload("Material").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_images.id ASC")
		}).
		Preload("Variants.Color").
		Preload("Variants.Size")
}

func (r *ProductRepositoryImpl) FindAll() ([]models.Product, error) {

	var products []models.Product

	err := preloadProduct(r.db).
		Order("display_order ASC, id DESC").
		Find(&products).
		Error

	return products, err
}

func (r *ProductRepositoryImpl) FindByID(
	id uint,
) (*models.Product, error) {

	var product models.Product

	err := preloadProduct(r.db).
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
	filter domainRepositories.ProductFilter,
) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	query := r.db.
		Model(&models.Product{}).
		Where("is_active = ?", true)

	if filter.Search != "" {
		normFA := utils.NormalizePersian(filter.Search)
		normEN := utils.NormalizeEnglish(filter.Search)
		raw := strings.TrimSpace(filter.Search)
		query = query.Where(
			"name_fa LIKE ? OR LOWER(name_en) LIKE ? OR product_code LIKE ?",
			"%"+normFA+"%",
			"%"+normEN+"%",
			"%"+raw+"%",
		)
	}

	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN (?)", filter.CategoryIDs)
	}

	if len(filter.MaterialIDs) > 0 {
		query = query.Where("material_id IN (?)", filter.MaterialIDs)
	}

	if len(filter.ColorIDs) > 0 {
		query = query.Where(
			"id IN (SELECT product_id FROM product_variants WHERE color_id IN (?))",
			filter.ColorIDs,
		)
	}

	if len(filter.SizeIDs) > 0 {
		query = query.Where(
			"id IN (SELECT product_id FROM product_variants WHERE size_id IN (?))",
			filter.SizeIDs,
		)
	}

	err := query.
		Count(&total).
		Error

	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	err = preloadProduct(query).
		Order("display_order ASC, id DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&products).
		Error

	return products, total, err
}

func (r *ProductRepositoryImpl) FindAdminProducts(
	filter domainRepositories.ProductFilter,
) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		normFA := utils.NormalizePersian(filter.Search)
		normEN := utils.NormalizeEnglish(filter.Search)
		raw := strings.TrimSpace(filter.Search)
		query = query.Where(
			"name_fa LIKE ? OR LOWER(name_en) LIKE ? OR product_code LIKE ?",
			"%"+normFA+"%",
			"%"+normEN+"%",
			"%"+raw+"%",
		)
	}

	err := query.
		Count(&total).
		Error

	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	err = preloadProduct(query).
		Order("display_order ASC, id DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&products).
		Error

	return products, total, err
}

func (r *ProductRepositoryImpl) FindActiveByID(
	id uint,
) (*models.Product, error) {

	var product models.Product

	err := preloadProduct(r.db).
		Where("is_active = ?", true).
		First(&product, id).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) Reorder(productIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range productIDs {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", id).
				Update("display_order", i+1).
				Error; err != nil {
				return err
			}
		}
		return nil
	})
}

