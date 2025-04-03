package services

import (
    "context"
    "errors"
    "log"

    "go.mongodb.org/mongo-driver/bson/primitive"
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
        GoogleID:       req.GoogleId,
        Name:           req.Name,
        Email:          req.Email,
        Password:       req.Password,
        ProfilePicture: req.ProfilePicture,
        Role:           "customer", 
    }

    err := s.UserUsecase.AddUser(user)
    if err != nil {
        log.Printf("Error adding user: %v", err)
        return nil, errors.New("failed to add user")
    }

    return &pb.AddUserResponse{
        User: &pb.User{
            Id:             user.ID.Hex(),
            GoogleId:       user.GoogleID,
            Name:           user.Name,
            Email:          user.Email,
            Password:       user.Password,
            ProfilePicture: user.ProfilePicture,
            Role:           user.Role,
        },
    }, nil
}


func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    objectID, err := primitive.ObjectIDFromHex(req.User.Id)
    if err != nil {
        log.Printf("Invalid ObjectID: %v", err)
        return nil, errors.New("invalid user ID")
    }

    updatedUser := &models.User{
        ID:             objectID,
        GoogleID:       req.User.GoogleId,
        Name:           req.User.Name,
        Email:          req.User.Email,
        Password:       req.User.Password,
        ProfilePicture: req.User.ProfilePicture,
        Role:           req.User.Role,
    }

    user, err := s.UserUsecase.UpdateUser(&req.User.Id, updatedUser)
    if err != nil {
        log.Printf("Error updating user: %v", err)
        return nil, errors.New("failed to update user")
    }

    return &pb.UpdateUserResponse{
        User: &pb.User{
            Id:             user.ID.Hex(),
            GoogleId:       user.GoogleID,
            Name:           user.Name,
            Email:          user.Email,
            Password:       user.Password,
            ProfilePicture: user.ProfilePicture,
            Role:           user.Role,
        },
    }, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.UserUsecase.GetUser(&req.Username)

    if err != nil {
        log.Printf("User not found: %v", req.Username)
        return nil, errors.New("user not found")
    }

    return &pb.GetUserResponse{
        User: &pb.User{
            Id:             user.ID.Hex(),
            GoogleId:       user.GoogleID,
            Name:           user.Name,
            Email:          user.Email,
            Password:       user.Password,
            ProfilePicture: user.ProfilePicture,
            Role:           user.Role,
        },
    }, nil
}


func (s *UserServiceImpl) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
    objectID, err := primitive.ObjectIDFromHex(req.Id)
    if err != nil {
        log.Printf("Invalid ObjectID: %v", err)
        return nil, errors.New("invalid user ID")
    }

    user, err := s.UserUsecase.GetUserByID(objectID.Hex())
    if err != nil {
        log.Printf("User not found: %v", err)
        return nil, errors.New("user not found")
    }

    return &pb.GetUserByIDResponse{
        User: &pb.User{
            Id:             user.ID.Hex(),
            GoogleId:       user.GoogleID,
            Name:           user.Name,
            Email:          user.Email,
            Password:       user.Password,
            ProfilePicture: user.ProfilePicture,
            Role:           user.Role,
        },
    }, nil
}


func (s *UserServiceImpl) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
    users, err := s.UserUsecase.GetAllUsers()
    if err != nil {
        log.Printf("Error retrieving users: %v", err)
        return nil, errors.New("failed to retrieve users")
    }

    var pbUsers []*pb.User
    for _, user := range users {
        pbUsers = append(pbUsers, &pb.User{
            Id:             user.ID.Hex(),
            GoogleId:       user.GoogleID,
            Name:           user.Name,
            Email:          user.Email,
            Password:       user.Password,
            ProfilePicture: user.ProfilePicture,
            Role:           user.Role,
        })
    }

    return &pb.GetAllUsersResponse{
        Users: pbUsers,
    }, nil
}


func (s *UserServiceImpl) PromoteUser(ctx context.Context, req *pb.PromoteUserRequest) (*pb.PromoteUserResponse, error) {
    err := s.UserUsecase.PromoteUser(req.Username)
    if err != nil {
        log.Printf("Error promoting user: %v", err)
        return nil, errors.New("failed to promote user")
    }

    return &pb.PromoteUserResponse{
        Message: "User promoted to admin successfully",
    }, nil
}