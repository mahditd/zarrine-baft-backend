package database

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"

	"gorm.io/gorm"
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