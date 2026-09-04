package services

import (
	"errors"
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type MaterialService struct {
	materialRepository repositories.MaterialRepository
}

func NewMaterialService(
	materialRepository repositories.MaterialRepository,
) *MaterialService {

	return &MaterialService{
		materialRepository: materialRepository,
	}
}


type CreateMaterialInput struct {
	NameFA string
	NameEN string
}


func (s *MaterialService) Create(
	input CreateMaterialInput,
) (*models.Material, error) {

	input.NameFA = normalizeMaterialName(input.NameFA)
	input.NameEN = normalizeMaterialName(input.NameEN)


	if input.NameFA == "" {
		return nil, errors.New("persian name is required")
	}

	if input.NameEN == "" {
		return nil, errors.New("english name is required")
	}


	existingFA, err := s.materialRepository.FindByNameFA(
		input.NameFA,
	)

	if err == nil && existingFA != nil {
		return nil, errors.New("persian material name already exists")
	}


	existingEN, err := s.materialRepository.FindByNameEN(
		input.NameEN,
	)

	if err == nil && existingEN != nil {
		return nil, errors.New("english material name already exists")
	}


	material := &models.Material{
		NameFA: input.NameFA,
		NameEN: input.NameEN,
	}


	err = s.materialRepository.Create(material)

	if err != nil {
		return nil, err
	}


	return material, nil
}


func (s *MaterialService) GetAll() ([]models.Material, error) {

	return s.materialRepository.FindAll()

}


func normalizeMaterialName(name string) string {

	name = strings.TrimSpace(name)

	name = strings.Join(
		strings.Fields(name),
		" ",
	)

	return name
}