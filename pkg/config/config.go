package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	MaxConns    int
	JwtSecret   string
	
	// Cấu hình Email
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Dùng 'db' thay vì 'localhost' để mặc định chạy được trong Docker
		dbURL = "postgres://postgres:postgrespassword@db:5432/tro_go?sslmode=disable"
	}

	maxConnsStr := os.Getenv("DB_MAX_CONNS")
	maxConns := 10
	if maxConnsStr != "" {
		if parsed, err := strconv.Atoi(maxConnsStr); err == nil {
			maxConns = parsed
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-super-secret-key-change-it-in-production"
	}

	smtpPort := 587 // Mặc định cho Gmail (TLS)
	if p := os.Getenv("SMTP_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			smtpPort = parsed
		}
	}

	return &Config{
		AppPort:      appPort,
		DatabaseURL:  dbURL,
		MaxConns:     maxConns,
		JwtSecret:    jwtSecret,
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}, nil
}
