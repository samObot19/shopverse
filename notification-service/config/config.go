package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}