package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type SizeRepositoryImpl struct {
	db *gorm.DB
}

func NewSizeRepository(
	db *gorm.DB,
) domainRepositories.SizeRepository {

	return &SizeRepositoryImpl{
		db: db,
	}
}

func (r *SizeRepositoryImpl) FindByID(
	id uint,
) (*models.Size, error) {

	var size models.Size

	err := r.db.
		First(&size, id).
		Error

	if err != nil {
		return nil, err
	}

	return &size, nil
}
