package services

import (
	"errors"
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type ProductService struct {
	productRepository  repositories.ProductRepository
	categoryRepository repositories.CategoryRepository
	materialRepository repositories.MaterialRepository
}

func NewProductService(
	productRepository repositories.ProductRepository,
	categoryRepository repositories.CategoryRepository,
	materialRepository repositories.MaterialRepository,
) *ProductService {

	return &ProductService{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		materialRepository: materialRepository,
	}
}

type CreateProductInput struct {
	NameFA string `json:"name_fa"`
	NameEN string `json:"name_en"`

	CategoryID uint `json:"category_id"`
	MaterialID uint `json:"material_id"`
}

func (s *ProductService) Create(
	input CreateProductInput,
) (*models.Product, error) {

	input.NameFA = strings.TrimSpace(input.NameFA)
	input.NameEN = strings.TrimSpace(input.NameEN)

	if input.NameFA == "" {
		return nil, errors.New("persian name is required")
	}

	if input.NameEN == "" {
		return nil, errors.New("english name is required")
	}

	_, err := s.categoryRepository.FindByID(
		input.CategoryID,
	)

	if err != nil {
		return nil, errors.New("category not found")
	}

	_, err = s.materialRepository.FindByID(
		input.MaterialID,
	)

	if err != nil {
		return nil, errors.New("material not found")
	}

	product := &models.Product{
		NameFA: input.NameFA,
		NameEN: input.NameEN,

		CategoryID: input.CategoryID,
		MaterialID: input.MaterialID,
	}

	err = s.productRepository.Create(product)

	if err != nil {
		return nil, err
	}

	createdProduct, err := s.productRepository.FindByID(
		product.ID,
	)

	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *ProductService) GetAll() ([]models.Product, error) {

	return s.productRepository.FindAll()

}
