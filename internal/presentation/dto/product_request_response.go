package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ProductRequestResponse struct {
	ID           uint                         `json:"id"`
	CustomerName string                       `json:"customer_name"`
	Phone        string                       `json:"phone"`
	Description  string                       `json:"description"`
	Status       string                       `json:"status"`
	Items        []ProductRequestItemResponse `json:"items"`
}

type ProductRequestItemResponse struct {
	ID uint `json:"id"`

	ProductVariantID uint `json:"product_variant_id"`

	ProductName string `json:"product_name"`
	ColorName   string `json:"color_name"`
	SizeName    string `json:"size_name"`

	Quantity int `json:"quantity"`
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
		Description:  request.Description,
		Status:       request.Status,
		Items:        items,
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
