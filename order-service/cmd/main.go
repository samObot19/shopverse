package main

import (
    "log"
    "net"

    "github.com/samObot19/shopverse/order-service/internal/config"
    "github.com/samObot19/shopverse/order-service/internal/events/subscribe"
    "github.com/samObot19/shopverse/order-service/internal/repository"
    "github.com/samObot19/shopverse/order-service/internal/services"
    "github.com/samObot19/shopverse/order-service/internal/usecases"
    productpb "github.com/samObot19/shopverse/order-service/clients/product-client/proto/pb"
    orderpb "github.com/samObot19/shopverse/order-service/proto/pb"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    conf := config.LoadConfig()
    db, err := config.InitDB(config.LoadDBConfig())
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    orderRepo := repository.NewOrderRepository(db)

    conn, err := grpc.Dial("localhost:" + conf.ProductPORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect to product service: %v", err)
    }
    defer conn.Close()
    productClient := productpb.NewProductServiceClient(conn)

    orderUsecase := usecases.NewOrderUsecase(orderRepo, productClient)
    orderService := services.NewOrderServiceServer(orderUsecase)

    go func() {
        log.Println("Starting stockEvent subscriber...")
        subscribe.SubscribeAndProcessStockEvent(orderUsecase)
    }()

    grpcServer := grpc.NewServer()
    orderpb.RegisterOrderServiceServer(grpcServer, orderService)

    listener, err := net.Listen("tcp", ":" + conf.OrderPORT)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    log.Println("gRPC server is running on port" + conf.OrderPORT)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}