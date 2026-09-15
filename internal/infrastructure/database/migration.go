package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

func Migrate(db *gorm.DB) {

	// Legacy databases created before the ProductCode feature have a
	// products table without product_code (or with NULLs). Adding the
	// NOT NULL column directly would abort with SQLSTATE 23502, so
	// backfill sequential 001-999 codes first. No data is deleted.
	if err := backfillProductCodes(db); err != nil {
		panic(fmt.Sprintf("database migration failed: %v", err))
	}

	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Material{},
		&models.Color{},
		&models.Size{},
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductImage{},
		&models.ProductRequest{},
		&models.ProductRequestItem{},
		&models.ProductRequestStatusHistory{},
	)

	if err != nil {
		panic(fmt.Sprintf("database migration failed: %v", err))
	}

	fmt.Println("Database migration completed")
}

// backfillProductCodes assigns unique 001-999 codes to legacy product rows
// that predate the product_code column, so AutoMigrate can safely apply
// the NOT NULL + unique constraints afterwards.
func backfillProductCodes(db *gorm.DB) error {
	if !db.Migrator().HasTable(&models.Product{}) {
		return nil // fresh database, nothing to backfill
	}

	if !db.Migrator().HasColumn(&models.Product{}, "product_code") {
		if err := db.Exec(`ALTER TABLE products ADD COLUMN product_code TEXT`).Error; err != nil {
			return err
		}
	}

	var nullIDs []uint
	if err := db.Table("products").
		Where("product_code IS NULL").
		Order("id").
		Pluck("id", &nullIDs).Error; err != nil {
		return err
	}

	if len(nullIDs) == 0 {
		return nil
	}

	var existing []string
	if err := db.Table("products").
		Where("product_code IS NOT NULL").
		Pluck("product_code", &existing).Error; err != nil {
		return err
	}

	used := make(map[string]bool, len(existing))
	for _, c := range existing {
		used[c] = true
	}

	next := 1
	for _, id := range nullIDs {
		for next <= 999 && used[fmt.Sprintf("%03d", next)] {
			next++
		}
		if next > 999 {
			return fmt.Errorf(
				"cannot backfill product_code: more than 999 products need codes (product id=%d)",
				id,
			)
		}

		code := fmt.Sprintf("%03d", next)
		used[code] = true
		next++

		if err := db.Exec(
			`UPDATE products SET product_code = ? WHERE id = ?`,
			code, id,
		).Error; err != nil {
			return err
		}

		fmt.Printf("Backfilled product id=%d with product_code=%s\n", id, code)
	}

	return nil
}
