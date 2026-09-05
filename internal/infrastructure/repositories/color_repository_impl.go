package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type ColorRepositoryImpl struct {
	db *gorm.DB
}

func NewColorRepository(
	db *gorm.DB,
) domainRepositories.ColorRepository {

	return &ColorRepositoryImpl{
		db: db,
	}
}

func (r *ColorRepositoryImpl) Create(
	color *models.Color,
) error {

	return r.db.Create(color).Error
}

func (r *ColorRepositoryImpl) FindByNameFA(
	name string,
) (*models.Color, error) {

	var color models.Color

	err := r.db.
		Where("name_fa = ?", name).
		First(&color).
		Error

	if err != nil {
		return nil, err
	}

	return &color, nil
}

func (r *ColorRepositoryImpl) FindByNameEN(
	name string,
) (*models.Color, error) {

	var color models.Color

	err := r.db.
		Where("name_en = ?", name).
		First(&color).
		Error

	if err != nil {
		return nil, err
	}

	return &color, nil
}

func (r *ColorRepositoryImpl) FindAll() ([]models.Color, error) {

	var colors []models.Color

	err := r.db.
		Find(&colors).
		Error

	return colors, err
}

func (r *ColorRepositoryImpl) FindByID(
	id uint,
) (*models.Color, error) {

	var color models.Color

	err := r.db.
		First(&color, id).
		Error

	if err != nil {
		return nil, err
	}

	return &color, nil
}
