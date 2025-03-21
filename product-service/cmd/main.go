package main

import (
    "log"
    "net"
    "os"

    "github.com/joho/godotenv"
    "google.golang.org/grpc"

    pb "github.com/samObot19/shopverse/product-service/proto/pb"
    "github.com/samObot19/shopverse/product-service/db"
    "github.com/samObot19/shopverse/product-service/repository"
    "github.com/samObot19/shopverse/product-service/service"
    "github.com/samObot19/shopverse/product-service/usecases"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Get environment variables
    mongoURI := os.Getenv("MONGO_URI")
    dbName := os.Getenv("DB_NAME")
    collectionName := os.Getenv("COLLECTION_NAME")
    grpcPort := os.Getenv("GRPC_PORT")

    if mongoURI == "" || dbName == "" || collectionName == "" || grpcPort == "" {
        log.Fatalf("Environment variables MONGO_URI, DB_NAME, COLLECTION_NAME, and GRPC_PORT must be set")
    }

    // Connect to MongoDB
    collection, err := db.NewMongoConnection(mongoURI, dbName, collectionName)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Initialize repository, use case, and gRPC server
    productRepo := repository.NewMongoProductRepository(collection)
    productUseCase := usecases.NewProductUseCase(productRepo)
    productServiceServer := service.NewProductServiceServer(productUseCase)

    // Start gRPC server
    listener, err := net.Listen("tcp", ":"+grpcPort)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterProductServiceServer(grpcServer, productServiceServer)

    log.Printf("gRPC server is running on port %s", grpcPort)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}