package dto

import (
	"time"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
)

type DashboardResponse struct {
	TotalProducts  int64                      `json:"total_products"`
	ActiveProducts int64                      `json:"active_products"`
	NewRequests    int64                      `json:"new_requests"`
	LatestRequests []DashboardRequestResponse `json:"latest_requests"`
}

type DashboardRequestResponse struct {
	ID            uint      `json:"id"`
	RequestNumber string    `json:"request_number"`
	CustomerName  string    `json:"customer_name"`
	CompanyName   string    `json:"company_name"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func FromDashboardResult(
	result *services.DashboardResult,
) DashboardResponse {

	latestRequests := make(
		[]DashboardRequestResponse,
		0,
		len(result.LatestRequests),
	)

	for _, request := range result.LatestRequests {

		latestRequests = append(
			latestRequests,
			DashboardRequestResponse{
				ID:            request.ID,
				RequestNumber: request.RequestNumber,
				CustomerName:  request.CustomerName,
				CompanyName:   request.CompanyName,
				Status:        string(request.Status),
				CreatedAt:     request.CreatedAt,
			},
		)
	}

	return DashboardResponse{
		TotalProducts:  result.ActiveProductsCount + result.InactiveProductsCount,
		ActiveProducts: result.ActiveProductsCount,
		NewRequests:    result.NewRequestsCount,
		LatestRequests: latestRequests,
	}
}
