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
	ClientURL      string

	AdminPhone    string
	AdminPassword string
	AdminName     string
}

func Load() *Config {
	// Load .env if present (ignore error if running with system env in Docker/CI)
	_ = godotenv.Load()

	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5173"
	}

	adminPhone := os.Getenv("ADMIN_PHONE")
	if adminPhone == "" {
		adminPhone = "09120000000"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin@123456"
	}

	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		adminName = "Super Admin"
	}

	return &Config{
		AppPort: appPort,

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: 24,
		UploadPath:     uploadPath,
		BaseURL:        os.Getenv("BASE_URL"),
		ClientURL:      clientURL,

		AdminPhone:    adminPhone,
		AdminPassword: adminPassword,
		AdminName:     adminName,
	}
}
