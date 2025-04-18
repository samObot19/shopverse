package main

import (
    "log"
    "net"
    "os"

    "github.com/joho/godotenv"
    pb "github.com/samObot19/shopverse/user-service/proto/pb"
    "github.com/samObot19/shopverse/user-service/repository"
    "github.com/samObot19/shopverse/user-service/services"
    "github.com/samObot19/shopverse/user-service/usecases"
    "google.golang.org/grpc"
)

func main() {
    if err := godotenv.Load(".env"); err != nil {
        log.Printf("No .env file found")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "50051" 
    }
    port = ":" + port

    listener, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", port, err)
    }

    userRepo := repository.NewMongoUserRepo()
    userUsecase := usecase.NewUserUsecase(userRepo)

    
    userService := &services.UserServiceImpl{
        UserUsecase: userUsecase,
    }

    grpcServer := grpc.NewServer()

    pb.RegisterUserServiceServer(grpcServer, userService)

    log.Printf("gRPC server is running on port %s", port)

    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}