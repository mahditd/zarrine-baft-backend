package services

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type SizeService struct {
	sizeRepository repositories.SizeRepository
}

func NewSizeService(
	sizeRepository repositories.SizeRepository,
) *SizeService {

	return &SizeService{
		sizeRepository: sizeRepository,
	}
}

func (s *SizeService) GetAll() ([]models.Size, error) {

	return s.sizeRepository.FindAll()
}
