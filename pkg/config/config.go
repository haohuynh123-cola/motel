package config

import (
	"log"
	"os"
	"strconv"
	"strings"

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

	// Cấu hình Kafka
	KafkaBrokers []string
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

	smtpPort := 587
	if p := os.Getenv("SMTP_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			smtpPort = parsed
		}
	}

	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if kafkaBrokersStr != "" {
		kafkaBrokers = strings.Split(kafkaBrokersStr, ",")
	} else {
		kafkaBrokers = []string{"kafka:9092"}
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
		KafkaBrokers: kafkaBrokers,
	}, nil
}
