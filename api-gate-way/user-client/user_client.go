package userclient

import (
	"context"
	"log"
	"time"

	pb "github.com/samObot19/shopverse/api-gate-way/user-client/proto/pb"
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


func ConnectToUserService(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
		return nil, err
	}
	return conn, nil
}


func (uc *UserClient) AddUser(ctx context.Context, name, email, password string, googleID, profilePicture string) (*pb.AddUserResponse, error) {
	req := &pb.AddUserRequest{
		Name:           name,
		Email:          email,
		Password:       password,
		GoogleId:       googleID,
		ProfilePicture: profilePicture,
	}
	resp, err := uc.client.AddUser(ctx, req)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		return nil, err
	}
	return resp, nil
}


func (uc *UserClient) UpdateUser(ctx context.Context, id, name, email, password, googleID, profilePicture, role string) (*pb.UpdateUserResponse, error) {
	req := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:             id,
			Name:           name,
			Email:          email,
			Password:       password,
			GoogleId:       googleID,
			ProfilePicture: profilePicture,
			Role:           role,
		},
	}
	resp, err := uc.client.UpdateUser(ctx, req)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, err
	}
	return resp, nil
}


func (uc *UserClient) GetUserByID(ctx context.Context, id string) (*pb.GetUserByIDResponse, error) {
	req := &pb.GetUserByIDRequest{
		Id: id,
	}
	resp, err := uc.client.GetUserByID(ctx, req)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)
		return nil, err
	}
	return resp, nil
}


func (uc *UserClient) GetAllUsers(ctx context.Context) (*pb.GetAllUsersResponse, error) {
	req := &pb.GetAllUsersRequest{}
	resp, err := uc.client.GetAllUsers(ctx, req)
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		return nil, err
	}
	return resp, nil
}


func (uc *UserClient) PromoteUser(ctx context.Context, username string) (*pb.PromoteUserResponse, error) {
	req := &pb.PromoteUserRequest{
		Username: username,
	}
	resp, err := uc.client.PromoteUser(ctx, req)
	if err != nil {
		log.Printf("Error promoting user: %v", err)
		return nil, err
	}
	return resp, nil
}

func (uc *UserClient) GetUser(ctx context.Context, username string) (*pb.GetUserResponse, error) {
	req := &pb.GetUserRequest{
		Username: username,
	}
	resp, err := uc.client.GetUser(ctx, req)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return resp, nil
}

