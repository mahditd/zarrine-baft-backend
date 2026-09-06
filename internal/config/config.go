package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
	UploadPath     string
	BaseURL        string
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: 24,
		UploadPath:     uploadPath,
		BaseURL:        os.Getenv("BASE_URL"),
	}
}
