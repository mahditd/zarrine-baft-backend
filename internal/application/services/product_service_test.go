package services_test

import (
	"errors"
	"testing"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type mockCategoryRepo struct {
	categories map[uint]*models.Category
}

func (m *mockCategoryRepo) Create(c *models.Category) error {
	c.ID = uint(len(m.categories) + 1)
	m.categories[c.ID] = c
	return nil
}
func (m *mockCategoryRepo) FindByNameFA(name string) (*models.Category, error) { return nil, nil }
func (m *mockCategoryRepo) FindByNameEN(name string) (*models.Category, error) { return nil, nil }
func (m *mockCategoryRepo) FindAll() ([]models.Category, error)                { return nil, nil }
func (m *mockCategoryRepo) FindByID(id uint) (*models.Category, error) {
	if c, ok := m.categories[id]; ok {
		return c, nil
	}
	return nil, errors.New("category not found")
}

type mockMaterialRepo struct {
	materials map[uint]*models.Material
}

func (m *mockMaterialRepo) Create(mat *models.Material) error {
	mat.ID = uint(len(m.materials) + 1)
	m.materials[mat.ID] = mat
	return nil
}
func (m *mockMaterialRepo) FindByNameFA(name string) (*models.Material, error) { return nil, nil }
func (m *mockMaterialRepo) FindByNameEN(name string) (*models.Material, error) { return nil, nil }
func (m *mockMaterialRepo) FindAll() ([]models.Material, error)                { return nil, nil }
func (m *mockMaterialRepo) FindByID(id uint) (*models.Material, error) {
	if mat, ok := m.materials[id]; ok {
		return mat, nil
	}
	return nil, errors.New("material not found")
}

type mockProductRepo struct {
	products map[uint]*models.Product
}

func (m *mockProductRepo) Create(p *models.Product) error {
	if p.ID == 0 {
		p.ID = uint(len(m.products) + 1)
	}
	// Shift existing
	for _, existing := range m.products {
		if existing.DisplayOrder >= 1 {
			existing.DisplayOrder++
		}
	}
	p.DisplayOrder = 1
	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepo) FindAll() ([]models.Product, error) {
	res := make([]models.Product, 0, len(m.products))
	for _, p := range m.products {
		res = append(res, *p)
	}
	return res, nil
}

func (m *mockProductRepo) FindByID(id uint) (*models.Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, errors.New("product not found")
}

func (m *mockProductRepo) FindByProductCode(code string) (*models.Product, error) {
	for _, p := range m.products {
		if p.ProductCode == code {
			return p, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockProductRepo) FindActiveProducts(filter repositories.ProductFilter) ([]models.Product, int64, error) {
	res := make([]models.Product, 0)
	for _, p := range m.products {
		if p.IsActive {
			res = append(res, *p)
		}
	}
	return res, int64(len(res)), nil
}

func (m *mockProductRepo) FindAdminProducts(filter repositories.ProductFilter) ([]models.Product, int64, error) {
	res := make([]models.Product, 0)
	for _, p := range m.products {
		if filter.IsActive == nil || p.IsActive == *filter.IsActive {
			res = append(res, *p)
		}
	}
	return res, int64(len(res)), nil
}

func (m *mockProductRepo) FindActiveByID(id uint) (*models.Product, error) {
	if p, ok := m.products[id]; ok && p.IsActive {
		return p, nil
	}
	return nil, errors.New("not found")
}

func (m *mockProductRepo) Update(product *models.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepo) Reorder(productIDs []uint) error {
	for i, id := range productIDs {
		if p, ok := m.products[id]; ok {
			p.DisplayOrder = i + 1
		}
	}
	return nil
}

func (m *mockProductRepo) Delete(product *models.Product) error {
	delete(m.products, product.ID)
	return nil
}

func (m *mockProductRepo) GetActiveCount() (int64, error) {
	return 0, nil
}

func (m *mockProductRepo) GetInactiveCount() (int64, error) {
	return 0, nil
}

type mockProductImageRepo struct {
	counts map[uint]int64
}

func (m *mockProductImageRepo) Create(image *models.ProductImage) error { return nil }
func (m *mockProductImageRepo) FindByID(id uint) (*models.ProductImage, error) {
	return nil, errors.New("not found")
}
func (m *mockProductImageRepo) FindByProductID(productID uint) ([]models.ProductImage, error) {
	return nil, nil
}
func (m *mockProductImageRepo) CountByProductID(productID uint) (int64, error) {
	if m.counts == nil {
		return 1, nil
	}
	if c, ok := m.counts[productID]; ok {
		return c, nil
	}
	return 1, nil
}
func (m *mockProductImageRepo) Update(image *models.ProductImage) error { return nil }
func (m *mockProductImageRepo) Reorder(productID uint, imageIDs []uint) error { return nil }
func (m *mockProductImageRepo) Delete(image *models.ProductImage) error { return nil }

func TestProductService_ProductCodeValidation(t *testing.T) {
	prodRepo := &mockProductRepo{products: make(map[uint]*models.Product)}
	catRepo := &mockCategoryRepo{categories: make(map[uint]*models.Category)}
	matRepo := &mockMaterialRepo{materials: make(map[uint]*models.Material)}

	_ = catRepo.Create(&models.Category{NameFA: "کت", NameEN: "Coat"})
	_ = matRepo.Create(&models.Material{NameFA: "پشم", NameEN: "Wool"})

	svc := services.NewProductService(prodRepo, catRepo, matRepo, &mockProductImageRepo{})

	// Valid 3-digit product code
	p1, err := svc.Create(services.CreateProductInput{
		ProductCode: "001",
		NameFA:      "كت   زمستاني",
		NameEN:      "Winter Coat",
		CategoryID:  1,
		MaterialID:  1,
	})
	if err != nil {
		t.Fatalf("expected valid creation, got error: %v", err)
	}
	if p1.ProductCode != "001" {
		t.Errorf("expected product code '001', got '%s'", p1.ProductCode)
	}
	// Verify Persian name was normalized: "كت   زمستاني" -> "کت زمستانی"
	if p1.NameFA != "کت زمستانی" {
		t.Errorf("expected normalized name 'کت زمستانی', got '%s'", p1.NameFA)
	}

	// Invalid codes: less than 3 digits, more than 3 digits, or "000"
	invalidCodes := []string{"1", "12", "1234", "abc", "000"}
	for _, code := range invalidCodes {
		_, err := svc.Create(services.CreateProductInput{
			ProductCode: code,
			NameFA:      "تست",
			NameEN:      "Test",
			CategoryID:  1,
			MaterialID:  1,
		})
		if err == nil {
			t.Errorf("expected error for invalid product code '%s', got nil", code)
		}
	}

	// Duplicate product code
	_, err = svc.Create(services.CreateProductInput{
		ProductCode: "001", // already used
		NameFA:      "دیگر کت",
		NameEN:      "Other Coat",
		CategoryID:  1,
		MaterialID:  1,
	})
	if err == nil || err.Error() != "product code already exists" {
		t.Errorf("expected 'product code already exists' error, got %v", err)
	}
}

func TestProductService_OrderingAndDeletionRules(t *testing.T) {
	prodRepo := &mockProductRepo{products: make(map[uint]*models.Product)}
	catRepo := &mockCategoryRepo{categories: make(map[uint]*models.Category)}
	matRepo := &mockMaterialRepo{materials: make(map[uint]*models.Material)}

	_ = catRepo.Create(&models.Category{NameFA: "کت", NameEN: "Coat"})
	_ = matRepo.Create(&models.Material{NameFA: "پشم", NameEN: "Wool"})

	svc := services.NewProductService(prodRepo, catRepo, matRepo, &mockProductImageRepo{})

	p1, _ := svc.Create(services.CreateProductInput{
		ProductCode: "001",
		NameFA:      "محصول اول",
		NameEN:      "Product 1",
		CategoryID:  1,
		MaterialID:  1,
	})
	p2, _ := svc.Create(services.CreateProductInput{
		ProductCode: "002",
		NameFA:      "محصول دوم",
		NameEN:      "Product 2",
		CategoryID:  1,
		MaterialID:  1,
	})

	// New products appear at the top: p2 should be display_order 1, p1 should be shifted to 2
	refreshedP2, _ := prodRepo.FindByID(p2.ID)
	refreshedP1, _ := prodRepo.FindByID(p1.ID)
	if refreshedP2.DisplayOrder != 1 {
		t.Errorf("expected new product p2 to have display_order 1, got %d", refreshedP2.DisplayOrder)
	}
	if refreshedP1.DisplayOrder != 2 {
		t.Errorf("expected product p1 to have display_order 2, got %d", refreshedP1.DisplayOrder)
	}

	// Reorder
	err := svc.Reorder([]uint{p1.ID, p2.ID})
	if err != nil {
		t.Fatalf("unexpected reorder error: %v", err)
	}
	refreshedP1, _ = prodRepo.FindByID(p1.ID)
	refreshedP2, _ = prodRepo.FindByID(p2.ID)
	if refreshedP1.DisplayOrder != 1 || refreshedP2.DisplayOrder != 2 {
		t.Errorf("expected p1 order 1, p2 order 2, got %d, %d", refreshedP1.DisplayOrder, refreshedP2.DisplayOrder)
	}

	// Deletion rejection
	err = svc.Delete(p1.ID)
	if err == nil || err.Error() != "products cannot be deleted; use deactivate instead" {
		t.Errorf("expected deletion rejection error, got %v", err)
	}
}
