package dto

import (
	"time"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductRequestResponse struct {
	ID            uint                                  `json:"id"`
	RequestNumber string                                `json:"request_number"`
	UserID        uint                                  `json:"user_id"`
	CustomerName  string                                `json:"customer_name"`
	Phone         string                                `json:"phone"`
	CompanyName   string                                `json:"company_name"`
	CompanyPhone  string                                `json:"company_phone"`
	Description   string                                `json:"description"`
	Status        string                                `json:"status"`
	Items         []ProductRequestItemResponse          `json:"items"`
	StatusHistory []ProductRequestStatusHistoryResponse `json:"status_history"`
	CreatedAt     time.Time                             `json:"created_at"`
	UpdatedAt     time.Time                             `json:"updated_at"`
}

type ProductRequestStatusHistoryResponse struct {
	ID         uint      `json:"id"`
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	Note       string    `json:"note,omitempty"`
	AdminName  string    `json:"admin_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProductRequestItemResponse struct {
	ID               uint      `json:"id"`
	ProductVariantID uint      `json:"product_variant_id"`
	ProductCode      string    `json:"product_code,omitempty"`
	ProductNameFA    string    `json:"product_name_fa,omitempty"`
	ProductNameEN    string    `json:"product_name_en,omitempty"`
	ColorNameFA      string    `json:"color_name_fa,omitempty"`
	ColorNameEN      string    `json:"color_name_en,omitempty"`
	SizeName         string    `json:"size_name,omitempty"`
	Quantity         int       `json:"quantity"`
	PriceSnapshot    int64     `json:"price_snapshot"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func mapItems(items []models.ProductRequestItem) []ProductRequestItemResponse {
	response := make([]ProductRequestItemResponse, 0, len(items))

	for _, item := range items {
		itemRes := ProductRequestItemResponse{
			ID:               item.ID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			PriceSnapshot:    item.PriceSnapshot,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		}

		if item.ProductVariant != nil {
			if item.ProductVariant.Product != nil {
				itemRes.ProductCode = item.ProductVariant.Product.ProductCode
				itemRes.ProductNameFA = item.ProductVariant.Product.NameFA
				itemRes.ProductNameEN = item.ProductVariant.Product.NameEN
			}

			if item.ProductVariant.Color != nil {
				itemRes.ColorNameFA = item.ProductVariant.Color.NameFA
				itemRes.ColorNameEN = item.ProductVariant.Color.NameEN
			}

			if item.ProductVariant.Size != nil {
				itemRes.SizeName = item.ProductVariant.Size.Name
			}
		}

		response = append(response, itemRes)
	}

	return response
}

func FromProductRequestForAdmin(request *models.ProductRequest) ProductRequestResponse {
	if request == nil {
		return ProductRequestResponse{}
	}

	history := make([]ProductRequestStatusHistoryResponse, 0, len(request.StatusHistory))
	for _, h := range request.StatusHistory {
		adminName := ""
		if h.Admin != nil {
			adminName = h.Admin.FullName
		}

		history = append(history, ProductRequestStatusHistoryResponse{
			ID:         h.ID,
			FromStatus: string(h.FromStatus),
			ToStatus:   string(h.ToStatus),
			Note:       h.Note,
			AdminName:  adminName,
			CreatedAt:  h.CreatedAt,
		})
	}

	return ProductRequestResponse{
		ID:            request.ID,
		RequestNumber: request.RequestNumber,
		UserID:        request.UserID,
		CustomerName:  request.CustomerName,
		Phone:         request.Phone,
		CompanyName:   request.CompanyName,
		CompanyPhone:  request.CompanyPhone,
		Description:   request.Description,
		Status:        string(request.Status),
		Items:         mapItems(request.Items),
		StatusHistory: history,
		CreatedAt:     request.CreatedAt,
		UpdatedAt:     request.UpdatedAt,
	}
}

func FromProductRequestForCustomer(request *models.ProductRequest) ProductRequestResponse {
	if request == nil {
		return ProductRequestResponse{}
	}

	history := make([]ProductRequestStatusHistoryResponse, 0, len(request.StatusHistory))
	for _, h := range request.StatusHistory {
		// As per SRS 13.2 & 14: Admin notes are private and never visible to customers.
		// Only customer cancellation notes are visible to the customer.
		note := ""
		if h.AdminID == nil && h.ToStatus == models.RequestCancelled {
			note = h.Note
		}

		history = append(history, ProductRequestStatusHistoryResponse{
			ID:         h.ID,
			FromStatus: string(h.FromStatus),
			ToStatus:   string(h.ToStatus),
			Note:       note,
			CreatedAt:  h.CreatedAt,
		})
	}

	return ProductRequestResponse{
		ID:            request.ID,
		RequestNumber: request.RequestNumber,
		UserID:        request.UserID,
		CustomerName:  request.CustomerName,
		Phone:         request.Phone,
		CompanyName:   request.CompanyName,
		CompanyPhone:  request.CompanyPhone,
		Description:   request.Description,
		Status:        string(request.Status),
		Items:         mapItems(request.Items),
		StatusHistory: history,
		CreatedAt:     request.CreatedAt,
		UpdatedAt:     request.UpdatedAt,
	}
}

// Aliases for compatibility
func FromProductRequest(request *models.ProductRequest) ProductRequestResponse {
	return FromProductRequestForAdmin(request)
}

func FromProductRequests(requests []models.ProductRequest) []ProductRequestResponse {
	response := make([]ProductRequestResponse, 0, len(requests))
	for _, request := range requests {
		response = append(response, FromProductRequestForAdmin(&request))
	}
	return response
}

func FromProductRequestsForCustomer(requests []models.ProductRequest) []ProductRequestResponse {
	response := make([]ProductRequestResponse, 0, len(requests))
	for _, request := range requests {
		response = append(response, FromProductRequestForCustomer(&request))
	}
	return response
}
