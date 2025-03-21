package db

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

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