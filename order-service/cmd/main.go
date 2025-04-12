package main

import (
    "log"
    "net"

    "github.com/samObot19/shopverse/order-service/internal/config"
    "github.com/samObot19/shopverse/order-service/internal/events/subscribe"
    "github.com/samObot19/shopverse/order-service/internal/repository"
    "github.com/samObot19/shopverse/order-service/internal/services"
    "github.com/samObot19/shopverse/order-service/internal/usecases"
    "github.com/samObot19/shopverse/order-service/proto/pb"

    "google.golang.org/grpc"
)

func main() {
    // Initialize database connection
    db, err := config.InitDB(config.LoadDBConfig())
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Initialize repository, usecase, and service
    orderRepo := repository.NewOrderRepository(db)
    orderUsecase := usecases.NewOrderUsecase(orderRepo)
    orderService := services.NewOrderServiceServer(orderUsecase)

    // Run the stockEvent subscriber in a separate goroutine
    go func() {
        log.Println("Starting stockEvent subscriber...")
        subscribe.SubscribeAndProcessStockEvent(orderUsecase)
    }()

    // Start gRPC server
    grpcServer := grpc.NewServer()
    pb.RegisterOrderServiceServer(grpcServer, orderService)

    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Println("gRPC server is running on port 50051")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}