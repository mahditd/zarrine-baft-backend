package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/utils"
)

func SeedSizes(db *gorm.DB) {

	sizes := []string{
		"S",
		"M",
		"L",
		"XL",
		"XXL",
		"XXXL",
		"XXXXL",
	}

	for _, size := range sizes {

		db.FirstOrCreate(
			&models.Size{},
			models.Size{
				Name: size,
			},
		)
	}
}

// SeedAdmin seeds the initial admin user if none exists (SRS 21).
func SeedAdmin(db *gorm.DB, phone, password, fullName string) error {
	var count int64
	err := db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Admin already exists
	}

	normalizedPhone, err := utils.NormalizePhone(phone)
	if err != nil {
		normalizedPhone = phone
	}

	// A previous test registration may already hold the seed phone.
	// Surface a clear error instead of a raw unique-constraint failure.
	var existing models.User
	err = db.Where("phone = ?", normalizedPhone).First(&existing).Error
	if err == nil {
		if existing.Role == models.RoleAdmin {
			return nil // Seed admin already exists
		}
		return fmt.Errorf(
			"seed phone %s already belongs to a non-admin user (id=%d); delete it or set a different ADMIN_PHONE: %w",
			normalizedPhone, existing.ID, errors.New("admin seed skipped"),
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		FullName:     fullName,
		Phone:        normalizedPhone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleAdmin,
	}

	return db.Create(&admin).Error
}
