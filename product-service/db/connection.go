package db

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql" // MySQL driver
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds the MySQL database configuration
type Config struct {
    MySQLUser     string
    MySQLPassword string
    MySQLHost     string
    MySQLPort     string
    MySQLDatabase string
    GRPCPort      string
}

// LoadConfig loads environment variables and returns the configuration
func LoadConfig() (*Config, error) {
    // Load environment variables from .env file
    err := godotenv.Load()
    if (err != nil) {
        log.Println("No .env file found, using system environment variables")
    }

    // Get environment variables
    config := &Config{
        MySQLUser:     os.Getenv("MYSQL_USER"),
        MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
        MySQLHost:     os.Getenv("MYSQL_HOST"),
        MySQLPort:     os.Getenv("MYSQL_PORT"),
        MySQLDatabase: os.Getenv("MYSQL_DATABASE"),
        GRPCPort:      os.Getenv("GRPC_PORT"),
    }

    // Validate required variables
    if config.MySQLUser == "" || config.MySQLPassword == "" || config.MySQLHost == "" || config.MySQLPort == "" || config.MySQLDatabase == "" || config.GRPCPort == "" {
        return nil, fmt.Errorf("missing required environment variables")
    }

    return config, nil
}

// NewMongoConnection establishes a new connection to the MongoDB database
func NewMongoConnection(uri string, dbName string, collectionName string) (*mongo.Collection, error) {
    // Set client options
    clientOptions := options.Client().ApplyURI(uri)

    // Connect to MongoDB
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
    }

    // Create a context with a timeout for the connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connect the client
    err = client.Connect(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Ping the database to verify the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    // Return the collection
    collection := client.Database(dbName).Collection(collectionName)
    return collection, nil
}

// NewSqlConnection establishes a new connection to the MySQL database
func NewSqlConnection(config *Config) (*sql.DB, error) {
    // Construct the MySQL Data Source Name (DSN)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        config.MySQLUser,
        config.MySQLPassword,
        config.MySQLHost,
        config.MySQLPort,
        config.MySQLDatabase,
    )

    // Open a connection to the MySQL database
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
    }

    // Test the connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = db.PingContext(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to ping MySQL: %w", err)
    }

    return db, nil
}

