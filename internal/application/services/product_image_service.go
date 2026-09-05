package services

import (
	"errors"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type ProductImageService struct {
	productImageRepository repositories.ProductImageRepository
	productRepository      repositories.ProductRepository
}

func NewProductImageService(
	productImageRepository repositories.ProductImageRepository,
	productRepository repositories.ProductRepository,
) *ProductImageService {

	return &ProductImageService{
		productImageRepository: productImageRepository,
		productRepository:      productRepository,
	}
}

type CreateProductImageInput struct {
	ProductID uint   `json:"product_id"`
	ImageURL  string `json:"image_url"`
}

func (s *ProductImageService) Create(
	input CreateProductImageInput,
) (*models.ProductImage, error) {

	_, err := s.productRepository.FindByID(
		input.ProductID,
	)

	if err != nil {
		return nil, errors.New("product not found")
	}

	if input.ImageURL == "" {
		return nil, errors.New("image url is required")
	}

	image := &models.ProductImage{
		ProductID: input.ProductID,
		ImageURL:  input.ImageURL,
	}

	err = s.productImageRepository.Create(
		image,
	)

	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *ProductImageService) GetByProductID(
	productID uint,
) ([]models.ProductImage, error) {

	return s.productImageRepository.FindByProductID(
		productID,
	)
}

func (s *ProductImageService) Delete(
	id uint,
) error {

	return s.productImageRepository.Delete(id)
}
