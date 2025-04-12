package main

import (
    "log"
    "net"

    "google.golang.org/grpc"

    pb "github.com/samObot19/shopverse/product-service/proto/pb"
    "github.com/samObot19/shopverse/product-service/db"
    "github.com/samObot19/shopverse/product-service/repository"
    "github.com/samObot19/shopverse/product-service/service"
    "github.com/samObot19/shopverse/product-service/usecases"
)

func main() {
    // Load configuration
    config, err := db.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Connect to MySQL
    sqlDB, err := db.NewSqlConnection(config)
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }
    defer sqlDB.Close()

    log.Println("Connected to MySQL successfully")

    // Initialize repository, use case, and gRPC server
    productRepo := repository.NewMySQLProductRepository(sqlDB)
    productUseCase := usecases.NewProductUseCase(productRepo)
    productServiceServer := service.NewProductServiceServer(productUseCase)

    // Start gRPC server
    listener, err := net.Listen("tcp", ":"+config.GRPCPort)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", config.GRPCPort, err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterProductServiceServer(grpcServer, productServiceServer)

    log.Printf("gRPC server is running on port %s", config.GRPCPort)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}