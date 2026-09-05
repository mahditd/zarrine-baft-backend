package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"time"
)

type ProductRequestResponse struct {
	ID           uint                         `json:"id"`
	CustomerName string                       `json:"customer_name"`
	Phone        string                       `json:"phone"`
	CompanyName  string                       `json:"company_name"`
	CompanyPhone string                       `json:"company_phone"`
	Description  string                       `json:"description"`
	Status       string                       `json:"status"`
	Items        []ProductRequestItemResponse `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductRequestItemResponse struct {
	ID uint `json:"id"`

	ProductVariantID uint `json:"product_variant_id"`

	ProductName string `json:"product_name"`
	ColorName   string `json:"color_name"`
	SizeName    string `json:"size_name"`

	Quantity int `json:"quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromProductRequest(
	request *models.ProductRequest,
) ProductRequestResponse {

	items := make([]ProductRequestItemResponse, 0)

	for _, item := range request.Items {

		response := ProductRequestItemResponse{
			ID:               item.ID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,

			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if item.ProductVariant != nil {

			if item.ProductVariant.Product != nil {
				response.ProductName = item.ProductVariant.Product.NameEN
			}

			if item.ProductVariant.Color != nil {
				response.ColorName = item.ProductVariant.Color.NameEN
			}

			if item.ProductVariant.Size != nil {
				response.SizeName = item.ProductVariant.Size.Name
			}
		}

		items = append(items, response)
	}

	return ProductRequestResponse{
		ID:           request.ID,
		CustomerName: request.CustomerName,
		Phone:        request.Phone,
		CompanyName:  request.CompanyName,
		CompanyPhone: request.CompanyPhone,
		Description:  request.Description,
		Status:       string(request.Status),
		Items:        items,

		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
	}
}

func FromProductRequests(
	requests []models.ProductRequest,
) []ProductRequestResponse {

	response := make(
		[]ProductRequestResponse,
		0,
	)

	for _, request := range requests {

		response = append(
			response,
			FromProductRequest(&request),
		)
	}

	return response
}
