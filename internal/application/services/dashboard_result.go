package services

import "github.com/mahditd/zarrine-baft-backend/internal/domain/models"

type DashboardResult struct {
	ActiveProductsCount   int64
	InactiveProductsCount int64
	NewRequestsCount      int64
	LatestRequests        []models.ProductRequest
}