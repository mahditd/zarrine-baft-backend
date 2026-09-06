package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/utils"
)

var productCodeRegex = regexp.MustCompile(`^[0-9]{3}$`)

type ProductService struct {
	productRepository      repositories.ProductRepository
	categoryRepository     repositories.CategoryRepository
	materialRepository     repositories.MaterialRepository
	productImageRepository repositories.ProductImageRepository
}

func NewProductService(
	productRepository repositories.ProductRepository,
	categoryRepository repositories.CategoryRepository,
	materialRepository repositories.MaterialRepository,
	productImageRepository repositories.ProductImageRepository,
) *ProductService {

	return &ProductService{
		productRepository:      productRepository,
		categoryRepository:     categoryRepository,
		materialRepository:     materialRepository,
		productImageRepository: productImageRepository,
	}
}

type CreateProductInput struct {
	ProductCode string `json:"product_code"`
	NameFA      string `json:"name_fa"`
	NameEN      string `json:"name_en"`

	CategoryID uint `json:"category_id"`
	MaterialID uint `json:"material_id"`
}

type UpdateProductInput struct {
	NameFA string `json:"name_fa"`
	NameEN string `json:"name_en"`

	CategoryID uint `json:"category_id"`
	MaterialID uint `json:"material_id"`
}

func (s *ProductService) Create(
	input CreateProductInput,
) (*models.Product, error) {

	input.ProductCode = strings.TrimSpace(input.ProductCode)
	input.NameFA = utils.NormalizePersian(input.NameFA)
	input.NameEN = strings.TrimSpace(input.NameEN)

	// SRS 4.2: Product code must be exactly 3 digits (001 - 999)
	if !productCodeRegex.MatchString(input.ProductCode) || input.ProductCode == "000" {
		return nil, errors.New("product code must be exactly 3 digits between 001 and 999")
	}

	existingCode, err := s.productRepository.FindByProductCode(input.ProductCode)
	if err == nil && existingCode != nil {
		return nil, errors.New("product code already exists")
	}

	if input.NameFA == "" {
		return nil, errors.New("persian name is required")
	}

	if input.NameEN == "" {
		return nil, errors.New("english name is required")
	}

	_, err = s.categoryRepository.FindByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	_, err = s.materialRepository.FindByID(input.MaterialID)
	if err != nil {
		return nil, errors.New("material not found")
	}

	product := &models.Product{
		ProductCode: input.ProductCode,
		NameFA:      input.NameFA,
		NameEN:      input.NameEN,
		CategoryID:  input.CategoryID,
		MaterialID:  input.MaterialID,
		IsActive:    true,
	}

	err = s.productRepository.Create(product)
	if err != nil {
		return nil, err
	}

	return s.productRepository.FindByID(product.ID)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.productRepository.FindAll()
}

func (s *ProductService) GetByID(
	id uint,
) (*models.Product, error) {
	return s.productRepository.FindByID(id)
}

func (s *ProductService) Update(
	id uint,
	input UpdateProductInput,
) (*models.Product, error) {

	product, err := s.productRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	input.NameFA = utils.NormalizePersian(input.NameFA)
	input.NameEN = strings.TrimSpace(input.NameEN)

	if input.NameFA == "" {
		return nil, errors.New("persian name is required")
	}

	if input.NameEN == "" {
		return nil, errors.New("english name is required")
	}

	_, err = s.categoryRepository.FindByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	_, err = s.materialRepository.FindByID(input.MaterialID)
	if err != nil {
		return nil, errors.New("material not found")
	}

	product.NameFA = input.NameFA
	product.NameEN = input.NameEN
	product.CategoryID = input.CategoryID
	product.MaterialID = input.MaterialID

	err = s.productRepository.Update(product)
	if err != nil {
		return nil, err
	}

	return s.productRepository.FindByID(product.ID)
}

func (s *ProductService) UpdateStatus(
	id uint,
	isActive bool,
) error {

	product, err := s.productRepository.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	// SRS 5.1: minimum 1 image per product. Block activation without images.
	if isActive && s.productImageRepository != nil {
		count, err := s.productImageRepository.CountByProductID(id)
		if err != nil {
			return err
		}
		if count < 1 {
			return errors.New("product must have at least one image to be activated")
		}
	}

	product.IsActive = isActive
	return s.productRepository.Update(product)
}

func (s *ProductService) Delete(
	id uint,
) error {
	// SRS Section 6: Products cannot be deleted.
	return errors.New("products cannot be deleted; use deactivate instead")
}

func (s *ProductService) Reorder(productIDs []uint) error {
	if len(productIDs) == 0 {
		return errors.New("product IDs list cannot be empty")
	}
	return s.productRepository.Reorder(productIDs)
}

func (s *ProductService) GetActiveProducts(
	filter repositories.ProductFilter,
) ([]models.Product, int64, error) {
	return s.productRepository.FindActiveProducts(filter)
}

func (s *ProductService) GetAdminProducts(
	filter repositories.ProductFilter,
) ([]models.Product, int64, error) {
	return s.productRepository.FindAdminProducts(filter)
}

func (s *ProductService) GetActiveByID(
	id uint,
) (*models.Product, error) {
	return s.productRepository.FindActiveByID(id)
}

