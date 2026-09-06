package services_test

import (
	"fmt"
	"testing"
	"time"
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type mockUserRepo struct {
	users map[uint]*models.User
}

func (m *mockUserRepo) Create(user *models.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByPhone(phone string) (*models.User, error) {
	for _, u := range m.users {
		if u.Phone == phone {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) FindByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email != nil && *u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) FindByID(id uint) (*models.User, error) {
	if u, exists := m.users[id]; exists {
		return u, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) Update(user *models.User) error {
	m.users[user.ID] = user
	return nil
}

type mockVariantRepo struct {
	variants map[uint]*models.ProductVariant
}

func (m *mockVariantRepo) Create(v *models.ProductVariant) error {
	if v.ID == 0 {
		v.ID = uint(len(m.variants) + 1)
	}
	m.variants[v.ID] = v
	return nil
}

func (m *mockVariantRepo) FindByID(id uint) (*models.ProductVariant, error) {
	if v, exists := m.variants[id]; exists {
		return v, nil
	}
	return nil, fmt.Errorf("variant not found")
}

func (m *mockVariantRepo) FindByProductID(productID uint) ([]models.ProductVariant, error) {
	return nil, nil
}

func (m *mockVariantRepo) FindAll() ([]models.ProductVariant, error) {
	return nil, nil
}

func (m *mockVariantRepo) FindByProductColorAndSize(productID, colorID, sizeID uint) (*models.ProductVariant, error) {
	return nil, nil
}

type mockRequestRepo struct {
	requests        map[uint]*models.ProductRequest
	statusHistories []models.ProductRequestStatusHistory
}

func (m *mockRequestRepo) Create(r *models.ProductRequest) error {
	r.ID = uint(len(m.requests) + 1)
	m.requests[r.ID] = r
	return nil
}

func (m *mockRequestRepo) Update(r *models.ProductRequest) error {
	m.requests[r.ID] = r
	return nil
}

func (m *mockVariantRepo) Delete(
	variant *models.ProductVariant,
) error {
	return nil
}

func (m *mockRequestRepo) FindLatest(
	limit int,
) ([]models.ProductRequest, error) {

	return []models.ProductRequest{}, nil
}

func (m *mockRequestRepo) GetNewRequestsCount() (int64, error) {

	return 0, nil
}

func (m *mockVariantRepo) Update(
	variant *models.ProductVariant,
) error {
	return nil
}

func (m *mockRequestRepo) FindByID(id uint) (*models.ProductRequest, error) {
	if r, exists := m.requests[id]; exists {
		return r, nil
	}
	return nil, fmt.Errorf("request not found")
}

func (m *mockRequestRepo) FindAll() ([]models.ProductRequest, error) {
	res := make([]models.ProductRequest, 0, len(m.requests))
	for _, r := range m.requests {
		res = append(res, *r)
	}
	return res, nil
}

func (m *mockRequestRepo) FindPaginated(
	page int,
	limit int,
	status string,
	search string,
) ([]models.ProductRequest, int64, error) {

	res := make([]models.ProductRequest, 0)

	for _, r := range m.requests {

		// status filter
		if status != "" && string(r.Status) != status {
			continue
		}

		// search filter
		if search != "" {
			search = strings.ToLower(search)

			matched :=
				strings.Contains(strings.ToLower(r.RequestNumber), search) ||
					strings.Contains(strings.ToLower(r.CustomerName), search) ||
					strings.Contains(strings.ToLower(r.CompanyName), search)

			if !matched {
				continue
			}
		}

		res = append(res, *r)
	}

	return res, int64(len(res)), nil
}

func (m *mockRequestRepo) FindByUserIDPaginated(userID uint, page, limit int) ([]models.ProductRequest, int64, error) {
	res := make([]models.ProductRequest, 0)
	for _, r := range m.requests {
		if r.UserID == userID {
			res = append(res, *r)
		}
	}
	return res, int64(len(res)), nil
}

func (m *mockRequestRepo) FindByIDAndUserID(id, userID uint) (*models.ProductRequest, error) {
	if r, exists := m.requests[id]; exists && r.UserID == userID {
		return r, nil
	}
	return nil, fmt.Errorf("request not found")
}

func (m *mockRequestRepo) CreateStatusHistory(h *models.ProductRequestStatusHistory) error {
	h.ID = uint(len(m.statusHistories) + 1)
	m.statusHistories = append(m.statusHistories, *h)
	return nil
}

func (m *mockRequestRepo) GetLatestRequestNumber(year int) (string, error) {
	var latest string
	for _, r := range m.requests {
		if len(r.RequestNumber) >= 4 && r.RequestNumber[:4] == fmt.Sprintf("%04d", year) {
			if r.RequestNumber > latest {
				latest = r.RequestNumber
			}
		}
	}
	return latest, nil
}

func setupTestService() (*services.ProductRequestService, *mockRequestRepo, *mockVariantRepo, *mockUserRepo) {
	reqRepo := &mockRequestRepo{requests: make(map[uint]*models.ProductRequest)}
	varRepo := &mockVariantRepo{variants: make(map[uint]*models.ProductVariant)}
	userRepo := &mockUserRepo{users: make(map[uint]*models.User)}

	svc := services.NewProductRequestService(reqRepo, varRepo, userRepo)
	return svc, reqRepo, varRepo, userRepo
}

func TestProductRequestCreation_SuccessWithItemMergingAndPriceSnapshot(t *testing.T) {
	svc, _, varRepo, userRepo := setupTestService()

	companyName := "Zarrine Corp"
	companyPhone := "02188888888"
	user := &models.User{
		ID:           1,
		FullName:     "Ali Reza",
		Phone:        "09121112233",
		Role:         models.RoleCustomer,
		CompanyName:  &companyName,
		CompanyPhone: &companyPhone,
	}
	_ = userRepo.Create(user)

	activeProduct := &models.Product{
		ID:          1,
		ProductCode: "001",
		NameFA:      "کت زمستانه",
		NameEN:      "Winter Coat",
		IsActive:    true,
	}

	variant1 := &models.ProductVariant{
		ID:        10,
		ProductID: 1,
		Product:   activeProduct,
		Price:     850000,
	}
	variant2 := &models.ProductVariant{
		ID:        20,
		ProductID: 1,
		Product:   activeProduct,
		Price:     920000,
	}
	_ = varRepo.Create(variant1)
	_ = varRepo.Create(variant2)

	input := services.CreateProductRequestInput{
		CompanyName:  "Zarrine Corp",
		CompanyPhone: "02188888888",
		Description:  "Urgent wholesale request",
		Items: []services.CreateProductRequestItemInput{
			{ProductVariantID: 10, Quantity: 5},
			{ProductVariantID: 20, Quantity: 2},
			{ProductVariantID: 10, Quantity: 3}, // duplicate variant 10 -> should merge to 8
		},
	}

	req, err := svc.Create(user.ID, input)
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}

	if req.Status != models.RequestNew {
		t.Errorf("expected status 'new', got '%s'", req.Status)
	}

	currentYear := time.Now().Year()
	expectedReqNum := fmt.Sprintf("%04d-000001", currentYear)
	if req.RequestNumber != expectedReqNum {
		t.Errorf("expected request number '%s', got '%s'", expectedReqNum, req.RequestNumber)
	}

	if len(req.Items) != 2 {
		t.Fatalf("expected 2 merged items, got %d", len(req.Items))
	}

	if req.Items[0].ProductVariantID != 10 || req.Items[0].Quantity != 8 {
		t.Errorf("expected merged variant 10 with qty 8, got qty %d", req.Items[0].Quantity)
	}

	if req.Items[0].PriceSnapshot != 850000 {
		t.Errorf("expected price snapshot 850000, got %d", req.Items[0].PriceSnapshot)
	}

	if req.Items[1].ProductVariantID != 20 || req.Items[1].Quantity != 2 {
		t.Errorf("expected variant 20 with qty 2, got qty %d", req.Items[1].Quantity)
	}

	if req.Items[1].PriceSnapshot != 920000 {
		t.Errorf("expected price snapshot 920000, got %d", req.Items[1].PriceSnapshot)
	}
}

func TestProductRequestCreation_CompanyRequired(t *testing.T) {
	svc, _, varRepo, userRepo := setupTestService()

	user := &models.User{
		ID:       1,
		FullName: "Sara",
		Phone:    "09129998877",
		Role:     models.RoleCustomer,
	}
	_ = userRepo.Create(user)

	activeProduct := &models.Product{ID: 1, IsActive: true}
	variant := &models.ProductVariant{ID: 1, Product: activeProduct, Price: 100000}
	_ = varRepo.Create(variant)

	// Missing company name
	input := services.CreateProductRequestInput{
		CompanyName:  "",
		CompanyPhone: "02112345678",
		Items:        []services.CreateProductRequestItemInput{{ProductVariantID: 1, Quantity: 1}},
	}
	_, err := svc.Create(user.ID, input)
	if err == nil || err.Error() != "company name is required" {
		t.Errorf("expected 'company name is required' error, got %v", err)
	}

	// Missing company phone
	input2 := services.CreateProductRequestInput{
		CompanyName:  "Test Co",
		CompanyPhone: "",
		Items:        []services.CreateProductRequestItemInput{{ProductVariantID: 1, Quantity: 1}},
	}
	_, err2 := svc.Create(user.ID, input2)
	if err2 == nil || err2.Error() != "company phone is required" {
		t.Errorf("expected 'company phone is required' error, got %v", err2)
	}
}

func TestProductRequestCreation_InactiveProductRejected(t *testing.T) {
	svc, _, varRepo, userRepo := setupTestService()

	user := &models.User{ID: 1, FullName: "Sara", Phone: "09129998877"}
	_ = userRepo.Create(user)

	inactiveProduct := &models.Product{ID: 2, NameFA: "پیراهن غیرفعال", IsActive: false}
	variant := &models.ProductVariant{ID: 5, Product: inactiveProduct, Price: 200000}
	_ = varRepo.Create(variant)

	input := services.CreateProductRequestInput{
		CompanyName:  "My Company",
		CompanyPhone: "02199999999",
		Items:        []services.CreateProductRequestItemInput{{ProductVariantID: 5, Quantity: 10}},
	}

	_, err := svc.Create(user.ID, input)
	if err == nil {
		t.Fatal("expected error for inactive product, got nil")
	}
}

func TestProductRequest_CustomerCancel(t *testing.T) {
	svc, reqRepo, varRepo, userRepo := setupTestService()

	user := &models.User{ID: 1, FullName: "Customer"}
	_ = userRepo.Create(user)

	product := &models.Product{ID: 1, IsActive: true}
	variant := &models.ProductVariant{ID: 1, Product: product, Price: 50000}
	_ = varRepo.Create(variant)

	input := services.CreateProductRequestInput{
		CompanyName:  "Company",
		CompanyPhone: "02122222222",
		Items:        []services.CreateProductRequestItemInput{{ProductVariantID: 1, Quantity: 1}},
	}

	req, err := svc.Create(user.ID, input)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	// Customer cancels new request
	err = svc.CancelByCustomer(req.ID, user.ID, "Changed my mind")
	if err != nil {
		t.Fatalf("expected successful cancellation, got: %v", err)
	}

	updatedReq, _ := reqRepo.FindByID(req.ID)
	if updatedReq.Status != models.RequestCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", updatedReq.Status)
	}

	// Trying to cancel again should fail
	err = svc.CancelByCustomer(req.ID, user.ID, "Cancel again")
	if err == nil {
		t.Errorf("expected error cancelling already cancelled request, got nil")
	}
}

func TestProductRequest_AdminUpdateStatus(t *testing.T) {
	svc, reqRepo, varRepo, userRepo := setupTestService()

	user := &models.User{ID: 1, FullName: "Customer"}
	_ = userRepo.Create(user)

	admin := &models.User{ID: 2, FullName: "Admin User", Role: models.RoleAdmin}
	_ = userRepo.Create(admin)

	product := &models.Product{ID: 1, IsActive: true}
	variant := &models.ProductVariant{ID: 1, Product: product, Price: 50000}
	_ = varRepo.Create(variant)

	input := services.CreateProductRequestInput{
		CompanyName:  "Company",
		CompanyPhone: "02122222222",
		Items:        []services.CreateProductRequestItemInput{{ProductVariantID: 1, Quantity: 1}},
	}

	req, _ := svc.Create(user.ID, input)

	// Admin updates status to contacted with note
	err := svc.UpdateStatus(req.ID, admin.ID, models.RequestContacted, "Called customer, discussed terms")
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}

	updatedReq, _ := reqRepo.FindByID(req.ID)
	if updatedReq.Status != models.RequestContacted {
		t.Errorf("expected status 'contacted', got '%s'", updatedReq.Status)
	}

	if len(reqRepo.statusHistories) < 2 {
		t.Fatalf("expected at least 2 status history entries, got %d", len(reqRepo.statusHistories))
	}

	latestHistory := reqRepo.statusHistories[len(reqRepo.statusHistories)-1]
	if latestHistory.ToStatus != models.RequestContacted {
		t.Errorf("expected ToStatus 'contacted', got '%s'", latestHistory.ToStatus)
	}
	if latestHistory.AdminID == nil || *latestHistory.AdminID != admin.ID {
		t.Errorf("expected AdminID 2, got %v", latestHistory.AdminID)
	}
	if latestHistory.Note != "Called customer, discussed terms" {
		t.Errorf("expected admin note, got '%s'", latestHistory.Note)
	}
}
