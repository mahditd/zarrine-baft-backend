package services

import (
	"errors"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type ProductVariantService struct {
	productVariantRepository repositories.ProductVariantRepository
	productRepository        repositories.ProductRepository
	colorRepository          repositories.ColorRepository
}

func NewProductVariantService(
	productVariantRepository repositories.ProductVariantRepository,
	productRepository repositories.ProductRepository,
	colorRepository repositories.ColorRepository,
) *ProductVariantService {

	return &ProductVariantService{
		productVariantRepository: productVariantRepository,
		productRepository:        productRepository,
		colorRepository:          colorRepository,
	}
}

type CreateProductVariantInput struct {
	ProductID uint `json:"product_id"`
	ColorID   uint `json:"color_id"`
}

func (s *ProductVariantService) Create(
	input CreateProductVariantInput,
) (*models.ProductVariant, error) {

	_, err := s.productRepository.FindByID(
		input.ProductID,
	)

	if err != nil {
		return nil, errors.New("product not found")
	}

	_, err = s.colorRepository.FindByID(
		input.ColorID,
	)

	if err != nil {
		return nil, errors.New("color not found")
	}

	existing, err := s.productVariantRepository.FindByProductAndColor(
		input.ProductID,
		input.ColorID,
	)

	if err == nil && existing != nil {
		return nil, errors.New("this color already exists for this product")
	}

	variant := &models.ProductVariant{
		ProductID: input.ProductID,
		ColorID:   input.ColorID,
	}

	err = s.productVariantRepository.Create(
		variant,
	)

	if err != nil {
		return nil, err
	}

	createdVariant, err := s.productVariantRepository.FindByID(
		variant.ID,
	)

	if err != nil {
		return nil, err
	}

	return createdVariant, nil
}

func (s *ProductVariantService) GetByProductID(
	productID uint,
) ([]models.ProductVariant, error) {

	return s.productVariantRepository.FindByProductID(
		productID,
	)
}
