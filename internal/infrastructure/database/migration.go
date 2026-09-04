package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Database migration completed")
}