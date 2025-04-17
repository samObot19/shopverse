package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	ProductPORT  string
	OrderPORT  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname"),
		JWTSecret:    getEnv("JWT_SECRET", "your_jwt_secret"),
		ProductPORT:  getEnv("PRODUCT_PORT", "8081"),
		OrderPORT:    getEnv("ORDER_PORT", "8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}