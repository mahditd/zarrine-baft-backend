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
		Status:       models.RequestPending,
	}

	for _, item := range input.Items {

		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than zero")
		}

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

func (s *ProductRequestService) GetAll() (
	[]models.ProductRequest,
	error,
) {

	return s.productRequestRepository.FindAll()
}

func (s *ProductRequestService) UpdateStatus(
	id uint,
	status models.ProductRequestStatus,
) error {

	if !isValidRequestStatus(string(status)) {
		return errors.New("invalid request status")
	}
	request, err := s.productRequestRepository.FindByID(id)

	if err != nil {
		return err
	}

	request.Status = status

	return s.productRequestRepository.Update(request)
}

func (s *ProductRequestService) GetByID(
	id uint,
) (*models.ProductRequest, error) {

	return s.productRequestRepository.FindByID(id)
}

func isValidRequestStatus(
	status string,
) bool {

	switch models.ProductRequestStatus(status) {

	case models.RequestPending,
		models.RequestReviewing,
		models.RequestApproved,
		models.RequestRejected,
		models.RequestCompleted:

		return true
	}

	return false
}

func (s *ProductRequestService) GetPaginated(
	page int,
	limit int,
	status string,
) (
	[]models.ProductRequest,
	int64,
	error,
) {

	return s.productRequestRepository.FindPaginated(
		page,
		limit,
		status,
	)
}
