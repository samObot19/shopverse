package user_client

import (
    "context"
    //"log"
    //"time"

    pb "github.com/samObot19/shopverse/api-gate-way/user-client/proto/pb" // Import the generated proto package
    "google.golang.org/grpc"
)

type UserClient struct {
    client pb.UserServiceClient
}


func NewUserClient(conn *grpc.ClientConn) *UserClient {
    return &UserClient{
        client: pb.NewUserServiceClient(conn),
    }
}

// AddUser calls the AddUser RPC
func (uc *UserClient) AddUser(ctx context.Context, name, email, password string) (*pb.AddUserResponse, error) {
    req := &pb.AddUserRequest{
        Name:     name,
        Email:    email,
        Password: password,
    }
    resp, err := uc.client.AddUser(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// GetUser calls the GetUser RPC
func (uc *UserClient) GetUser(ctx context.Context, username string) (*pb.GetUserResponse, error) {
    req := &pb.GetUserRequest{Username: username}
    resp, err := uc.client.GetUser(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// PromoteUser calls the PromoteUser RPC
func (uc *UserClient) PromoteUser(ctx context.Context, username string) (*pb.PromoteUserResponse, error) {
    req := &pb.PromoteUserRequest{Username: username}
    resp, err := uc.client.PromoteUser(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// GetAllUsers calls the GetAllUsers RPC
func (uc *UserClient) GetAllUsers(ctx context.Context) (*pb.GetAllUsersResponse, error) {
    req := &pb.GetAllUsersRequest{}
    resp, err := uc.client.GetAllUsers(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// Example usage
// func main() {
//     // Connect to the gRPC server
//     conn, err := grpc.Dial("localhost:50500", grpc.WithInsecure()) 
//     if err != nil {
//         log.Fatalf("Failed to connect to server: %v", err)
//     }
//     defer conn.Close()

//     // Create a new UserClient
//     client := NewUserClient(conn)

//     // Example: Add a user
//     ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//     defer cancel()

//     addResp, err := client.AddUser(ctx, "samuel f/s", "sam@example.com", "password123")
//     if err != nil {
//         log.Fatalf("Error adding user: %v", err)
//     }
//     log.Printf("User added: %v", addResp.User)

//     // Example: Get a user
//     getResp, err := client.GetUser(ctx, "John Doe")
//     if err != nil {
//         log.Fatalf("Error fetching user: %v", err)
//     }
//     log.Printf("User fetched: %v", getResp.User)

    
//     promoteResp, err := client.PromoteUser(ctx, "John Doe")
//     if err != nil {
//         log.Fatalf("Error promoting user: %v", err)
//     }
//     log.Printf("User promoted: %s", promoteResp.Message)

//     // Example: Get all users
//     allUsersResp, err := client.GetAllUsers(ctx)
//     if err != nil {
//         log.Fatalf("Error fetching all users: %v", err)
//     }
//     log.Printf("All users: %v", allUsersResp.Users)
// }


