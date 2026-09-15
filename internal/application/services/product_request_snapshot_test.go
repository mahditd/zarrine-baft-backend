package services_test

import (
	"testing"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

func TestProductRequestCreation_SnapshotAllFields(t *testing.T) {
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

	product := &models.Product{
		ID:          1,
		ProductCode: "042",
		NameFA:      "کت زمستانه",
		NameEN:      "Winter Coat",
		IsActive:    true,
	}
	color := &models.Color{ID: 7, NameFA: "قرمز", NameEN: "Red"}
	size := &models.Size{ID: 3, Name: "XL"}
	variant := &models.ProductVariant{
		ID:        10,
		ProductID: 1,
		Product:   product,
		ColorID:   7,
		Color:     color,
		SizeID:    3,
		Size:      size,
		Price:     850000,
	}
	_ = varRepo.Create(variant)

	req, err := svc.Create(user.ID, services.CreateProductRequestInput{
		CompanyName:  "Zarrine Corp",
		CompanyPhone: "02188888888",
		Items:        []services.CreateProductRequestItemInput{{ProductVariantID: 10, Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}
	if len(req.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(req.Items))
	}
	item := req.Items[0]
	if item.ProductCode != "042" {
		t.Errorf("ProductCode: expected '042', got '%s'", item.ProductCode)
	}
	if item.ProductNameFA != "کت زمستانه" {
		t.Errorf("ProductNameFA: got '%s'", item.ProductNameFA)
	}
	if item.ProductNameEN != "Winter Coat" {
		t.Errorf("ProductNameEN: got '%s'", item.ProductNameEN)
	}
	if item.ColorNameFA != "قرمز" {
		t.Errorf("ColorNameFA: got '%s'", item.ColorNameFA)
	}
	if item.ColorNameEN != "Red" {
		t.Errorf("ColorNameEN: got '%s'", item.ColorNameEN)
	}
	if item.SizeName != "XL" {
		t.Errorf("SizeName: got '%s'", item.SizeName)
	}
	if item.PriceSnapshot != 850000 {
		t.Errorf("PriceSnapshot: expected 850000, got %d", item.PriceSnapshot)
	}
}

func TestProductRequestCreation_BlockedWhenColorOrSizeUnavailable(t *testing.T) {
	setup := func() (*services.ProductRequestService, *mockUserRepo) {
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
		product := &models.Product{ID: 1, ProductCode: "042", NameFA: "کت", NameEN: "Coat", IsActive: true}
		// Color soft-deleted => preload yields nil Color
		_ = varRepo.Create(&models.ProductVariant{
			ID: 10, ProductID: 1, Product: product,
			ColorID: 7, Color: nil,
			SizeID: 3, Size: &models.Size{ID: 3, Name: "XL"},
			Price: 100,
		})
		// Size soft-deleted => preload yields nil Size
		_ = varRepo.Create(&models.ProductVariant{
			ID: 11, ProductID: 1, Product: product,
			ColorID: 7, Color: &models.Color{ID: 7, NameFA: "قرمز", NameEN: "Red"},
			SizeID: 4, Size: nil,
			Price: 100,
		})
		return svc, userRepo
	}

	svc, userRepo := setup()
	user, _ := userRepo.FindByID(1)
	_, err := svc.Create(user.ID, services.CreateProductRequestInput{
		CompanyName: "Zarrine Corp", CompanyPhone: "02188888888",
		Items: []services.CreateProductRequestItemInput{{ProductVariantID: 10, Quantity: 1}},
	})
	if err == nil {
		t.Errorf("expected block when color unavailable, got nil")
	}

	svc, userRepo = setup()
	user, _ = userRepo.FindByID(1)
	_, err = svc.Create(user.ID, services.CreateProductRequestInput{
		CompanyName: "Zarrine Corp", CompanyPhone: "02188888888",
		Items: []services.CreateProductRequestItemInput{{ProductVariantID: 11, Quantity: 1}},
	})
	if err == nil {
		t.Errorf("expected block when size unavailable, got nil")
	}
}
