package main

import (
    "log"

    "github.com/samObot19/shopverse/order-service/internal/config"
    //"order-service/internal/events"
    "github.com/samObot19/shopverse/order-service/internal/repository"
    "github.com/samObot19/shopverse/order-service/internal/services"
    "github.com/samObot19/shopverse/order-service/internal/usecases"
    "github.com/samObot19/shopverse/order-service/proto/pb"

    "google.golang.org/grpc"
    "net"
)

func main() {
    // Load environment variables
    // kafkaBroker := os.Getenv("KAFKA_BROKER")
    // kafkaTopic := os.Getenv("KAFKA_TOPIC")

    // // Initialize Kafka publisher
    // kafkaPublisher, err := events.NewKafkaPublisher(kafkaBroker, kafkaTopic)
    // if err != nil {
    //     log.Fatalf("Failed to initialize Kafka publisher: %v", err)
    // }
    // defer kafkaPublisher.Close()

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