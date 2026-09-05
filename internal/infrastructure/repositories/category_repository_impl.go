package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(
	db *gorm.DB,
) domainRepositories.CategoryRepository {

	return &CategoryRepositoryImpl{
		db: db,
	}
}

func (r *CategoryRepositoryImpl) Create(
	category *models.Category,
) error {

	return r.db.Create(category).Error
}

func (r *CategoryRepositoryImpl) FindByNameFA(
	name string,
) (*models.Category, error) {

	var category models.Category

	err := r.db.
		Where("name_fa = ?", name).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepositoryImpl) FindByNameEN(
	name string,
) (*models.Category, error) {

	var category models.Category

	err := r.db.
		Where("name_en = ?", name).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepositoryImpl) FindAll() ([]models.Category, error) {

	var categories []models.Category

	err := r.db.
		Find(&categories).
		Error

	return categories, err
}

func (r *CategoryRepositoryImpl) FindByID(
	id uint,
) (*models.Category, error) {

	var category models.Category

	err := r.db.
		First(&category, id).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}
