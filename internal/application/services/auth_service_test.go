package services_test

import (
	"testing"

	"github.com/mahditd/zarrine-baft-backend/internal/application/services"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

func TestAuthService_GetAndUpdateProfile(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[uint]*models.User)}
	authService := services.NewAuthService(userRepo, "secret", 24)

	initialEmail := "user1@example.com"
	user1 := &models.User{
		ID:       1,
		FullName: "Old Name",
		Phone:    "09121112233",
		Email:    &initialEmail,
		Role:     models.RoleCustomer,
	}
	_ = userRepo.Create(user1)

	otherEmail := "user2@example.com"
	user2 := &models.User{
		ID:       2,
		FullName: "Other User",
		Phone:    "09124445566",
		Email:    &otherEmail,
		Role:     models.RoleCustomer,
	}
	_ = userRepo.Create(user2)

	// Test GetProfile
	profile, err := authService.GetProfile(user1.ID)
	if err != nil {
		t.Fatalf("expected to get profile, got error: %v", err)
	}
	if profile.FullName != "Old Name" {
		t.Errorf("expected 'Old Name', got '%s'", profile.FullName)
	}

	// Test UpdateProfile successfully
	updateInput := services.UpdateProfileInput{
		FullName:     "New Name",
		Email:        "user1_new@example.com",
		CompanyName:  "Brand New Co",
		CompanyPhone: "+982188889999",
		Country:      "Iran",
		Address:      "Tehran, Valiasr St",
	}

	updated, err := authService.UpdateProfile(user1.ID, updateInput)
	if err != nil {
		t.Fatalf("expected successful update, got error: %v", err)
	}

	if updated.FullName != "New Name" {
		t.Errorf("expected 'New Name', got '%s'", updated.FullName)
	}
	if updated.Email == nil || *updated.Email != "user1_new@example.com" {
		t.Errorf("expected updated email, got %v", updated.Email)
	}
	if updated.CompanyName == nil || *updated.CompanyName != "Brand New Co" {
		t.Errorf("expected 'Brand New Co', got %v", updated.CompanyName)
	}
	if updated.CompanyPhone == nil || *updated.CompanyPhone != "02188889999" {
		t.Errorf("expected normalized phone '02188889999', got %v", updated.CompanyPhone)
	}
	if updated.Country == nil || *updated.Country != "Iran" {
		t.Errorf("expected 'Iran', got %v", updated.Country)
	}
	if updated.Address == nil || *updated.Address != "Tehran, Valiasr St" {
		t.Errorf("expected 'Tehran, Valiasr St', got %v", updated.Address)
	}
	// Phone must stay intact
	if updated.Phone != "09121112233" {
		t.Errorf("login phone should remain '09121112233', got '%s'", updated.Phone)
	}

	// Test duplicate email error
	conflictInput := services.UpdateProfileInput{
		Email: "user2@example.com", // owned by user 2
	}
	_, err = authService.UpdateProfile(user1.ID, conflictInput)
	if err == nil || err.Error() != "email already exists" {
		t.Errorf("expected 'email already exists' error, got %v", err)
	}
}
