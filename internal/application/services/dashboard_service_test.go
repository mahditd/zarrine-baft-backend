package services_test

import (
	"testing"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type mockDashboardProductRepo struct {
	activeCount   int64
	inactiveCount int64
}

func (m *mockDashboardProductRepo) GetActiveCount() (int64, error) {
	return m.activeCount, nil
}

func (m *mockDashboardProductRepo) Create(
	product *models.Product,
) error {
	return nil
}

func (m *mockDashboardProductRepo) Update(
	product *models.Product,
) error {
	return nil
}

func (m *mockDashboardProductRepo) Delete(
	product *models.Product,
) error {
	return nil
}

func (m *mockDashboardProductRepo) FindActiveByID(
	id uint,
) (*models.Product, error) {
	return nil, nil
}

func (m *mockDashboardProductRepo) FindActiveProducts(
	filter repositories.ProductFilter,
) ([]models.Product, int64, error) {
	return nil, 0, nil
}

func (m *mockDashboardProductRepo) FindByID(
	id uint,
) (*models.Product, error) {
	return nil, nil
}

func (m *mockDashboardProductRepo) FindAdminProducts(
	filter repositories.ProductFilter,
) ([]models.Product, int64, error) {
	return nil, 0, nil
}

func (m *mockDashboardProductRepo) FindByProductCode(
	code string,
) (*models.Product, error) {
	return nil, nil
}

func (m *mockDashboardProductRepo) FindAll() (
	[]models.Product,
	error,
) {
	return nil, nil
}

func (m *mockDashboardProductRepo) Reorder(
	productIDs []uint,
) error {
	return nil
}

func (m *mockDashboardProductRepo) GetInactiveCount() (int64, error) {
	return m.inactiveCount, nil
}

type mockDashboardRequestRepo struct {
	newCount int64
	latest   []models.ProductRequest
}

func (m *mockDashboardRequestRepo) Create(
	request *models.ProductRequest,
) error {
	return nil
}

func (m *mockDashboardRequestRepo) Update(
	request *models.ProductRequest,
) error {
	return nil
}

func (m *mockDashboardRequestRepo) FindByID(
	id uint,
) (*models.ProductRequest, error) {
	return nil, nil
}

func (m *mockDashboardRequestRepo) FindAll() (
	[]models.ProductRequest,
	error,
) {
	return nil, nil
}

func (m *mockDashboardRequestRepo) FindPaginated(
	page int,
	limit int,
	status string,
	userRole string,
) ([]models.ProductRequest, int64, error) {
	return nil, 0, nil
}

func (m *mockDashboardRequestRepo) FindByUserIDPaginated(
	userID uint,
	page int,
	limit int,
) ([]models.ProductRequest, int64, error) {
	return nil, 0, nil
}

func (m *mockDashboardRequestRepo) FindByIDAndUserID(
	id uint,
	userID uint,
) (*models.ProductRequest, error) {
	return nil, nil
}

func (m *mockDashboardRequestRepo) CreateStatusHistory(
	history *models.ProductRequestStatusHistory,
) error {
	return nil
}

func (m *mockDashboardRequestRepo) GetLatestRequestNumber(
	year int,
) (string, error) {
	return "", nil
}

func (m *mockDashboardRequestRepo) GetNewRequestsCount() (int64, error) {
	return m.newCount, nil
}

func (m *mockDashboardRequestRepo) FindLatest(
	limit int,
) ([]models.ProductRequest, error) {

	if limit > len(m.latest) {
		return m.latest, nil
	}

	return m.latest[:limit], nil
}

func TestDashboardService_GetDashboard_Success(t *testing.T) {

	productRepo := &mockDashboardProductRepo{
		activeCount:   8,
		inactiveCount: 2,
	}

	requestRepo := &mockDashboardRequestRepo{
		newCount: 3,
		latest: []models.ProductRequest{
			{
				ID:            1,
				RequestNumber: "2026-000001",
				CustomerName:  "Ali",
				CompanyName:   "Zarrine Corp",
				Status:        models.RequestNew,
			},
		},
	}

	service := services.NewDashboardService(
		productRepo,
		requestRepo,
	)

	result, err := service.GetDashboard()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ActiveProductsCount != 8 {
		t.Errorf(
			"expected active products 8, got %d",
			result.ActiveProductsCount,
		)
	}

	if result.InactiveProductsCount != 2 {
		t.Errorf(
			"expected inactive products 2, got %d",
			result.InactiveProductsCount,
		)
	}

	if result.NewRequestsCount != 3 {
		t.Errorf(
			"expected new requests 3, got %d",
			result.NewRequestsCount,
		)
	}

	if len(result.LatestRequests) != 1 {
		t.Fatalf(
			"expected 1 latest request, got %d",
			len(result.LatestRequests),
		)
	}

	if result.LatestRequests[0].RequestNumber != "2026-000001" {
		t.Errorf(
			"unexpected request number: %s",
			result.LatestRequests[0].RequestNumber,
		)
	}
}
