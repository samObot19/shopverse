package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql" 
)


type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}


func LoadDBConfig() *DBConfig {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	return &DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}
}

// InitDB initializes and returns a MySQL database connection
func InitDB(config *DBConfig) (*sql.DB, error) {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Name)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
