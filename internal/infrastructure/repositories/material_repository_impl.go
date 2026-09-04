package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type MaterialRepositoryImpl struct {
	db *gorm.DB
}

func NewMaterialRepository(
	db *gorm.DB,
) domainRepositories.MaterialRepository {

	return &MaterialRepositoryImpl{
		db: db,
	}
}


func (r *MaterialRepositoryImpl) Create(
	material *models.Material,
) error {

	return r.db.Create(material).Error
}


func (r *MaterialRepositoryImpl) FindByNameFA(
	name string,
) (*models.Material, error) {

	var material models.Material

	err := r.db.
		Where("name_fa = ?", name).
		First(&material).
		Error

	if err != nil {
		return nil, err
	}

	return &material, nil
}


func (r *MaterialRepositoryImpl) FindByNameEN(
	name string,
) (*models.Material, error) {

	var material models.Material

	err := r.db.
		Where("name_en = ?", name).
		First(&material).
		Error

	if err != nil {
		return nil, err
	}

	return &material, nil
}


func (r *MaterialRepositoryImpl) FindAll() ([]models.Material, error) {

	var materials []models.Material

	err := r.db.
		Find(&materials).
		Error

	return materials, err
}

func (r *MaterialRepositoryImpl) FindByID(
	id uint,
) (*models.Material, error) {

	var material models.Material

	err := r.db.
		First(&material, id).
		Error

	if err != nil {
		return nil, err
	}

	return &material, nil
}