package services

import (
	"errors"
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type CategoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(
	categoryRepository repositories.CategoryRepository,
) *CategoryService {

	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

type CreateCategoryInput struct {
	NameFA string
	NameEN string
}

func (s *CategoryService) Create(
	input CreateCategoryInput,
) (*models.Category, error) {

	input.NameFA = normalizeName(input.NameFA)
	input.NameEN = normalizeName(input.NameEN)

	if input.NameFA == "" {
		return nil, errors.New("persian name is required")
	}

	if input.NameEN == "" {
		return nil, errors.New("english name is required")
	}

	existingFA, err := s.categoryRepository.FindByNameFA(
		input.NameFA,
	)

	if err == nil && existingFA != nil {
		return nil, errors.New("persian category name already exists")
	}

	existingEN, err := s.categoryRepository.FindByNameEN(
		input.NameEN,
	)

	if err == nil && existingEN != nil {
		return nil, errors.New("english category name already exists")
	}

	category := &models.Category{
		NameFA: input.NameFA,
		NameEN: input.NameEN,
	}

	err = s.categoryRepository.Create(category)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func normalizeName(name string) string {

	name = strings.TrimSpace(name)

	name = strings.Join(
		strings.Fields(name),
		" ",
	)

	return name
}

func (s *CategoryService) GetAll() ([]models.Category, error) {

	return s.categoryRepository.FindAll()

}
