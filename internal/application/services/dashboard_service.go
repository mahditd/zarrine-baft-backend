package services

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type DashboardService struct {
	productRepository        repositories.ProductRepository
	productRequestRepository repositories.ProductRequestRepository
}

func NewDashboardService(
	productRepository repositories.ProductRepository,
	productRequestRepository repositories.ProductRequestRepository,
) *DashboardService {

	return &DashboardService{
		productRepository:        productRepository,
		productRequestRepository: productRequestRepository,
	}
}

func (s *DashboardService) GetDashboard() (
	*DashboardResult,
	error,
) {

	activeProducts, err := s.productRepository.GetActiveCount()

	if err != nil {
		return nil, err
	}

	inactiveProducts, err := s.productRepository.GetInactiveCount()

	if err != nil {
		return nil, err
	}

	newRequests, err := s.productRequestRepository.GetNewRequestsCount()

	if err != nil {
		return nil, err
	}

	latestRequests, err := s.productRequestRepository.FindLatest(10)

	if err != nil {
		return nil, err
	}

	return &DashboardResult{
		ActiveProductsCount:   activeProducts,
		InactiveProductsCount: inactiveProducts,
		NewRequestsCount:      newRequests,
		LatestRequests:        latestRequests,
	}, nil
}
