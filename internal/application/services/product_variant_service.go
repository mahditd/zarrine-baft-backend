package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type ProductVariantService struct {
	productVariantRepository repositories.ProductVariantRepository
	productRepository        repositories.ProductRepository
	colorRepository          repositories.ColorRepository
	sizeRepository           repositories.SizeRepository
}

func NewProductVariantService(
	productVariantRepository repositories.ProductVariantRepository,
	productRepository repositories.ProductRepository,
	colorRepository repositories.ColorRepository,
	sizeRepository repositories.SizeRepository,
) *ProductVariantService {

	return &ProductVariantService{
		productVariantRepository: productVariantRepository,
		productRepository:        productRepository,
		colorRepository:          colorRepository,
		sizeRepository:           sizeRepository,
	}
}

type CreateProductVariantInput struct {
	ProductID uint  `json:"product_id"`
	ColorID   uint  `json:"color_id"`
	SizeID    uint  `json:"size_id"`
	Price     int64 `json:"price"`
}

func (s *ProductVariantService) Create(
	input CreateProductVariantInput,
) (*models.ProductVariant, error) {
	if input.ProductID == 0 ||
		input.ColorID == 0 ||
		input.SizeID == 0 ||
		input.Price <= 0 {
		return nil, errors.New("invalid variant data")
	}

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

	_, err = s.sizeRepository.FindByID(
		input.SizeID,
	)

	if err != nil {
		return nil, errors.New("size not found")
	}

	existing, err := s.productVariantRepository.FindByProductColorAndSize(
		input.ProductID,
		input.ColorID,
		input.SizeID,
	)

	if err == nil && existing != nil {
		return nil, errors.New("this variant already exists for this product")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	variant := &models.ProductVariant{
		ProductID: input.ProductID,
		ColorID:   input.ColorID,
		SizeID:    input.SizeID,
		Price:     input.Price,
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
