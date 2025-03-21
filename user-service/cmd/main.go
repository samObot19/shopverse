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
    // Load environment variables from .env file
    if err := godotenv.Load(".env"); err != nil {
        log.Printf("No .env file found")
    }

    // Get the port from the environment variable, default to ":50051" if not set
    port := os.Getenv("PORT")
    if port == "" {
        port = "50051" // Default port
    }
    port = ":" + port

    // Create a listener on the specified port
    listener, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", port, err)
    }

    // Initialize the repository
    userRepo := repository.NewMongoUserRepo()

    // Initialize the usecase
    userUsecase := usecase.NewUserUsecase(userRepo)

    // Initialize the gRPC service implementation
    userService := &services.UserServiceImpl{
        UserUsecase: userUsecase,
    }

    // Create a new gRPC server
    grpcServer := grpc.NewServer()

    // Register the UserService with the gRPC server
    pb.RegisterUserServiceServer(grpcServer, userService)

    log.Printf("gRPC server is running on port %s", port)

    // Start the gRPC server
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}