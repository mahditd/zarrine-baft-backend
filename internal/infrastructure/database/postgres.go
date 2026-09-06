package database

import (
	"fmt"

	"github.com/mahditd/zarrine-baft-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {

	fmt.Println("DATABASE:", cfg.DBName)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var dbName string
	db.Raw("SELECT current_database()").Scan(&dbName)

	fmt.Println("CONNECTED DATABASE:", dbName)

	return db
}
