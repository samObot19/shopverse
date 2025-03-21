package services

import (
    "context"
    "errors"
    "log"

    "github.com/samObot19/shopverse/user-service/models"
    pb "github.com/samObot19/shopverse/user-service/proto/pb"
    "github.com/samObot19/shopverse/user-service/usecases"
)

type UserServiceImpl struct {
    pb.UnimplementedUserServiceServer
    UserUsecase *usecase.UserUsecase
}


func (s *UserServiceImpl) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     "customer", // Default role
    }

    // Call the usecase to add the user
    err := s.UserUsecase.AddUser(user)
    if err != nil {
        log.Printf("Error adding user: %v", err)
        return nil, errors.New("failed to add user")
    }

    // Return the created user
    return &pb.AddUserResponse{
        User: &pb.User{
            Id:       "", // MongoDB will generate the ID
            Name:     user.Name,
            Email:    user.Email,
            Password: user.Password,
            Role:     user.Role,
        },
    }, nil
}

// UpdateUser updates an existing user's details
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    updatedUser := &models.User{
        Name:     req.User.Name,
        Email:    req.User.Email,
        Password: req.User.Password,
        Role:     req.User.Role,
    }

    // Call the usecase to update the user
    user, err := s.UserUsecase.UpdateUser(&req.User.Id, updatedUser)
    if err != nil {
        log.Printf("Error updating user: %v", err)
        return nil, errors.New("failed to update user")
    }

    // Return the updated user
    return &pb.UpdateUserResponse{
        User: &pb.User{
            Id:       user.ID,
            Name:     user.Name,
            Email:    user.Email,
            Password: user.Password,
            Role:     user.Role,
        },
    }, nil
}

// PromoteUser promotes a user to admin
func (s *UserServiceImpl) PromoteUser(ctx context.Context, req *pb.PromoteUserRequest) (*pb.PromoteUserResponse, error) {
    // Call the usecase to promote the user
    err := s.UserUsecase.PromoteUser(req.Username)
    if err != nil {
        log.Printf("Error promoting user: %v", err)
        return nil, errors.New("failed to promote user")
    }

    return &pb.PromoteUserResponse{
        Message: "User promoted to admin successfully",
    }, nil
}

// GetUser retrieves a user by username
func (s *UserServiceImpl) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Call the usecase to get the user
    user, err := s.UserUsecase.GetUser(&req.Username)
    if err != nil {
        log.Printf("User not found: %s", req.Username)
        return nil, errors.New("user not found")
    }

    // Return the user
    return &pb.GetUserResponse{
        User: &pb.User{
            Id:       user.ID,
            Name:     user.Name,
            Email:    user.Email,
            Password: user.Password,
            Role:     user.Role,
        },
    }, nil
}

// GetAllUsers retrieves all users
func (s *UserServiceImpl) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
    // Call the usecase to get all users
    users, err := s.UserUsecase.GetAllUsers()
    if err != nil {
        log.Printf("Error retrieving users: %v", err)
        return nil, errors.New("failed to retrieve users")
    }

    // Convert users to protobuf format
    var pbUsers []*pb.User
    for _, user := range users {
        pbUsers = append(pbUsers, &pb.User{
            Id:       user.ID,
            Name:     user.Name,
            Email:    user.Email,
            Password: user.Password,
            Role:     user.Role,
        })
    }

    // Return the list of users
    return &pb.GetAllUsersResponse{
        Users: pbUsers,
    }, nil
}