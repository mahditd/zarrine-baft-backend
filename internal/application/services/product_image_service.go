package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/infrastructure/storage"
)

type ProductImageService struct {
	productImageRepository repositories.ProductImageRepository
	productRepository      repositories.ProductRepository
	storage                storage.Storage
	baseURL                string
}

func NewProductImageService(
	productImageRepository repositories.ProductImageRepository,
	productRepository repositories.ProductRepository,
	fileStorage storage.Storage,
	baseURL string,
) *ProductImageService {

	return &ProductImageService{
		productImageRepository: productImageRepository,
		productRepository:      productRepository,
		storage:                fileStorage,
		baseURL:                baseURL,
	}
}

type UploadProductImageInput struct {
	ProductID uint
	File      *multipart.FileHeader
}

func (s *ProductImageService) Upload(
	input UploadProductImageInput,
) (*models.ProductImage, error) {

	_, err := s.productRepository.FindByID(
		input.ProductID,
	)

	if err != nil {
		return nil, errors.New("product not found")
	}

	count, err := s.productImageRepository.CountByProductID(input.ProductID)

	if err != nil {
		return nil, err
	}

	if count >= 15 {
		return nil, errors.New("maximum 15 images allowed per product")
	}

	if input.File == nil {
		return nil, errors.New("image file is required")
	}

	if input.File.Size > 10*1024*1024 {
		return nil, errors.New("image size must be less than 10MB")
	}

	file, err := input.File.Open()

	if err != nil {
		return nil, err
	}

	defer file.Close()

	buffer := make([]byte, 512)

	_, err = file.Read(buffer)

	if err != nil {
		return nil, err
	}

	mimeType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[mimeType] {
		return nil, fmt.Errorf(
			"unsupported image type: %s",
			mimeType,
		)
	}

	path, err := s.storage.Save(
		input.File,
	)

	if err != nil {
		return nil, err
	}

	imageURL := "/" + filepath.ToSlash(path)

	if s.baseURL != "" {
		imageURL = s.baseURL + imageURL
	}

	image := &models.ProductImage{
		ProductID: input.ProductID,
		FilePath:  path,
		ImageURL:  imageURL,
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

	image, err := s.productImageRepository.FindByID(id)

	if err != nil {
		return errors.New("image not found")
	}

	count, err := s.productImageRepository.CountByProductID(
		image.ProductID,
	)

	if err != nil {
		return err
	}

	if count <= 1 {
		return errors.New("product must have at least one image")
	}

	if image.FilePath != "" {

		err := s.storage.Delete(
			image.FilePath,
		)

		if err != nil {
			return err
		}
	}

	return s.productImageRepository.Delete(image)
}

func (s *ProductImageService) Reorder(
	productID uint,
	imageIDs []uint,
) error {

	if len(imageIDs) == 0 {
		return errors.New("image ids are required")
	}

	if len(imageIDs) > 15 {
		return errors.New("maximum 15 images allowed")
	}

	return s.productImageRepository.Reorder(
		productID,
		imageIDs,
	)
}
