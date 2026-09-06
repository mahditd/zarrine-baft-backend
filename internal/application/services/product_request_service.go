package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
	"github.com/mahditd/zarrine-baft-backend/internal/utils"
)

type ProductRequestService struct {
	productRequestRepository repositories.ProductRequestRepository
	productVariantRepository repositories.ProductVariantRepository
	userRepository           repositories.UserRepository
}

func NewProductRequestService(
	productRequestRepository repositories.ProductRequestRepository,
	productVariantRepository repositories.ProductVariantRepository,
	userRepository repositories.UserRepository,
) *ProductRequestService {

	return &ProductRequestService{
		productRequestRepository: productRequestRepository,
		productVariantRepository: productVariantRepository,
		userRepository:           userRepository,
	}
}

type CreateProductRequestItemInput struct {
	ProductVariantID uint `json:"product_variant_id"`
	Quantity         int  `json:"quantity"`
}

type CreateProductRequestInput struct {
	CustomerName string `json:"customer_name"`
	Phone        string `json:"phone"`
	CompanyName  string `json:"company_name"`
	CompanyPhone string `json:"company_phone"`
	Description  string `json:"description"`

	Items []CreateProductRequestItemInput `json:"items"`
}

func (s *ProductRequestService) generateRequestNumber() (string, error) {
	year := time.Now().Year()

	latest, err := s.productRequestRepository.GetLatestRequestNumber(year)
	if err != nil {
		return "", err
	}

	nextSeq := 1
	if latest != "" {
		parts := strings.Split(latest, "-")
		if len(parts) == 2 {
			if seq, err := strconv.Atoi(parts[1]); err == nil {
				nextSeq = seq + 1
			}
		}
	}

	return fmt.Sprintf("%04d-%06d", year, nextSeq), nil
}

func (s *ProductRequestService) Create(
	userID uint,
	input CreateProductRequestInput,
) (*models.ProductRequest, error) {

	user, err := s.userRepository.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if len(input.Items) == 0 {
		return nil, errors.New("request must contain at least one item")
	}

	// Determine customer name & phone (fallback to profile if empty)
	customerName := strings.TrimSpace(input.CustomerName)
	if customerName == "" {
		customerName = user.FullName
	}

	customerPhone := strings.TrimSpace(input.Phone)
	if customerPhone == "" {
		customerPhone = user.Phone
	} else {
		norm, err := utils.NormalizePhone(customerPhone)
		if err != nil {
			return nil, errors.New("invalid phone number")
		}
		customerPhone = norm
	}

	// Company Name & Company Phone are strictly required
	companyName := strings.TrimSpace(input.CompanyName)
	if companyName == "" && user.CompanyName != nil {
		companyName = strings.TrimSpace(*user.CompanyName)
	}
	if companyName == "" {
		return nil, errors.New("company name is required")
	}

	companyPhone := strings.TrimSpace(input.CompanyPhone)
	if companyPhone == "" && user.CompanyPhone != nil {
		companyPhone = strings.TrimSpace(*user.CompanyPhone)
	}
	if companyPhone == "" {
		return nil, errors.New("company phone is required")
	}
	normCompanyPhone, err := utils.NormalizePhone(companyPhone)
	if err != nil {
		return nil, errors.New("invalid company phone number")
	}
	companyPhone = normCompanyPhone

	// Merge duplicate variants and validate quantities (SRS 11.2 & 11.3)
	mergedQuantities := make(map[uint]int)
	variantOrder := make([]uint, 0)

	for _, item := range input.Items {
		if item.ProductVariantID == 0 {
			return nil, errors.New("invalid product variant id")
		}

		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be at least 1")
		}

		if item.Quantity > 999 {
			return nil, errors.New("quantity cannot exceed 999")
		}

		if _, exists := mergedQuantities[item.ProductVariantID]; !exists {
			variantOrder = append(variantOrder, item.ProductVariantID)
		}
		mergedQuantities[item.ProductVariantID] += item.Quantity

		if mergedQuantities[item.ProductVariantID] > 999 {
			return nil, errors.New("merged quantity cannot exceed 999 for any item")
		}
	}

	// Validate variants exist, product is active, and capture price snapshot (SRS 12.1 & 12.2)
	requestItems := make([]models.ProductRequestItem, 0, len(variantOrder))

	for _, variantID := range variantOrder {
		variant, err := s.productVariantRepository.FindByID(variantID)
		if err != nil || variant == nil {
			return nil, errors.New("product variant not found")
		}

		if variant.Product != nil && !variant.Product.IsActive {
			return nil, fmt.Errorf("product '%s' is inactive and cannot be requested", variant.Product.NameFA)
		}

		qty := mergedQuantities[variantID]

		requestItems = append(
			requestItems,
			models.ProductRequestItem{
				ProductVariantID: variantID,
				Quantity:         qty,
				PriceSnapshot:    variant.Price,
			},
		)
	}

	reqNumber, err := s.generateRequestNumber()
	if err != nil {
		return nil, errors.New("failed to generate request number: " + err.Error())
	}

	request := &models.ProductRequest{
		UserID:        userID,
		RequestNumber: reqNumber,
		CustomerName:  customerName,
		Phone:         customerPhone,
		CompanyName:   companyName,
		CompanyPhone:  companyPhone,
		Description:   strings.TrimSpace(input.Description),
		Status:        models.RequestNew,
		Items:         requestItems,
	}

	err = s.productRequestRepository.Create(request)
	if err != nil {
		return nil, err
	}

	// Record initial status history
	initialHistory := models.ProductRequestStatusHistory{
		RequestID:  request.ID,
		FromStatus: "",
		ToStatus:   models.RequestNew,
		Note:       "Request submitted",
	}
	_ = s.productRequestRepository.CreateStatusHistory(&initialHistory)

	return s.productRequestRepository.FindByID(request.ID)
}

func (s *ProductRequestService) GetAll() (
	[]models.ProductRequest,
	error,
) {
	return s.productRequestRepository.FindAll()
}

func (s *ProductRequestService) GetByID(
	id uint,
) (*models.ProductRequest, error) {
	return s.productRequestRepository.FindByID(id)
}

func (s *ProductRequestService) GetPaginated(
	page int,
	limit int,
	status string,
	search string,
) (
	[]models.ProductRequest,
	int64,
	error,
) {
	return s.productRequestRepository.FindPaginated(
		page,
		limit,
		status,
		search,
	)
}

func (s *ProductRequestService) GetByUserIDPaginated(
	userID uint,
	page int,
	limit int,
) ([]models.ProductRequest, int64, error) {
	return s.productRequestRepository.FindByUserIDPaginated(
		userID,
		page,
		limit,
	)
}

func (s *ProductRequestService) GetByIDAndUserID(
	id uint,
	userID uint,
) (*models.ProductRequest, error) {
	return s.productRequestRepository.FindByIDAndUserID(id, userID)
}

func (s *ProductRequestService) UpdateStatus(
	id uint,
	adminID uint,
	status models.ProductRequestStatus,
	note string,
) error {

	if !isValidRequestStatus(string(status)) {
		return errors.New("invalid request status")
	}

	request, err := s.productRequestRepository.FindByID(id)
	if err != nil {
		return err
	}

	oldStatus := request.Status
	if oldStatus == status {
		return errors.New("status is already set to " + string(status))
	}

	request.Status = status

	err = s.productRequestRepository.Update(request)
	if err != nil {
		return err
	}

	history := models.ProductRequestStatusHistory{
		RequestID:  request.ID,
		FromStatus: oldStatus,
		ToStatus:   status,
		Note:       strings.TrimSpace(note),
		AdminID:    &adminID,
	}

	return s.productRequestRepository.CreateStatusHistory(&history)
}

func (s *ProductRequestService) CancelByCustomer(
	id uint,
	userID uint,
	reason string,
) error {

	request, err := s.productRequestRepository.FindByIDAndUserID(id, userID)
	if err != nil {
		return errors.New("request not found")
	}

	if request.Status != models.RequestNew {
		return errors.New("only requests with status 'new' can be cancelled")
	}

	oldStatus := request.Status
	request.Status = models.RequestCancelled

	err = s.productRequestRepository.Update(request)
	if err != nil {
		return err
	}

	note := strings.TrimSpace(reason)
	if note == "" {
		note = "Cancelled by customer"
	}

	history := models.ProductRequestStatusHistory{
		RequestID:  request.ID,
		FromStatus: oldStatus,
		ToStatus:   models.RequestCancelled,
		Note:       note,
		AdminID:    nil,
	}

	return s.productRequestRepository.CreateStatusHistory(&history)
}

func isValidRequestStatus(status string) bool {
	switch models.ProductRequestStatus(status) {
	case models.RequestNew,
		models.RequestContacted,
		models.RequestInDiscussion,
		models.RequestCompleted,
		models.RequestCancelled:
		return true
	}

	return false
}
