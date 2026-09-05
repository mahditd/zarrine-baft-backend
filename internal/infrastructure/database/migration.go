package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

func Migrate(db *gorm.DB) {

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
	)

	if err != nil {
		panic(fmt.Sprintf("database migration failed: %v", err))
	}

	fmt.Println("Database migration completed")
}
