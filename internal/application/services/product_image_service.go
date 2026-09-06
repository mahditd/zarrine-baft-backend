package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

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

func validateImageFile(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("image file is required")
	}

	if file.Size > 10*1024*1024 {
		return "", errors.New("image size must be less than 10MB")
	}

	// SRS 5.1: extension check (in addition to MIME check).
	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return "", errors.New("unsupported image extension: must be JPG, JPEG, PNG or WebP")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)
	switch mimeType {
	case "image/jpeg", "image/png", "image/webp":
		return mimeType, nil
	default:
		return "", fmt.Errorf("unsupported image type: %s", mimeType)
	}
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

	if _, err := validateImageFile(input.File); err != nil {
		return nil, err
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

type ReplaceProductImageInput struct {
	ProductID uint
	ImageID   uint
	File      *multipart.FileHeader
}

// Replace swaps the file of an existing image in place (SRS 5.3):
// position, display_order and cover flag are preserved.
func (s *ProductImageService) Replace(
	input ReplaceProductImageInput,
) (*models.ProductImage, error) {
	image, err := s.productImageRepository.FindByID(input.ImageID)
	if err != nil {
		return nil, errors.New("image not found")
	}

	if image.ProductID != input.ProductID {
		return nil, errors.New("image does not belong to this product")
	}

	if _, err := validateImageFile(input.File); err != nil {
		return nil, err
	}

	newPath, err := s.storage.Save(input.File)
	if err != nil {
		return nil, err
	}

	newURL := "/" + filepath.ToSlash(newPath)
	if s.baseURL != "" {
		newURL = s.baseURL + newURL
	}

	oldPath := image.FilePath
	image.FilePath = newPath
	image.ImageURL = newURL

	if err := s.productImageRepository.Update(image); err != nil {
		_ = s.storage.Delete(newPath)
		return nil, err
	}

	_ = s.storage.Delete(oldPath)

	return image, nil
}
