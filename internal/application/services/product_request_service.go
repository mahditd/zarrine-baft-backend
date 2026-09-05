package services

import (
	"errors"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type ProductRequestService struct {
	productRequestRepository repositories.ProductRequestRepository
	productVariantRepository repositories.ProductVariantRepository
}

func NewProductRequestService(
	productRequestRepository repositories.ProductRequestRepository,
	productVariantRepository repositories.ProductVariantRepository,
) *ProductRequestService {

	return &ProductRequestService{
		productRequestRepository: productRequestRepository,
		productVariantRepository: productVariantRepository,
	}
}

type CreateProductRequestItemInput struct {
	ProductVariantID uint `json:"product_variant_id"`
	Quantity         int  `json:"quantity"`
}

type CreateProductRequestInput struct {
	CustomerName string `json:"customer_name"`

	Phone string `json:"phone"`

	Description string `json:"description"`

	Items []CreateProductRequestItemInput `json:"items"`
}

func (s *ProductRequestService) Create(
	input CreateProductRequestInput,
) (*models.ProductRequest, error) {

	if len(input.Items) == 0 {
		return nil, errors.New("request must contain at least one item")
	}

	request := &models.ProductRequest{
		CustomerName: input.CustomerName,
		Phone:        input.Phone,
		Description:  input.Description,
		Status:       "pending",
	}

	for _, item := range input.Items {

		_, err := s.productVariantRepository.FindByID(
			item.ProductVariantID,
		)

		if err != nil {
			return nil, errors.New("product variant not found")
		}

		request.Items = append(
			request.Items,
			models.ProductRequestItem{
				ProductVariantID: item.ProductVariantID,
				Quantity:         item.Quantity,
			},
		)
	}

	err := s.productRequestRepository.Create(request)

	if err != nil {
		return nil, err
	}

	createdRequest, err := s.productRequestRepository.FindByID(
		request.ID,
	)

	if err != nil {
		return nil, err
	}

	return createdRequest, nil
}
