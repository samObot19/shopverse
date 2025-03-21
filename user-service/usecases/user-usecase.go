package usecase

import (
	"github.com/samObot19/shopverse/user-service/models"
	"github.com/samObot19/shopverse/user-service/repository"
	"errors"
)

type UserUsecase struct{
	db repository.UserRepository
}


func NewUserUsecase(con repository.UserRepository) *UserUsecase{
	return &UserUsecase{
		db : con,
	}
}

func (s *UserUsecase) GetUser(username *string) (models.User, error) {
	user, ok := s.db.ReadUser(*username)

	if !ok{
		return models.User{}, errors.New("the user not found")
	}
	return user, nil
}

func (s *UserUsecase) PromoteUser(username string) error{
	return s.db.ChangeRoleToAdmin(username)
}

func (s *UserUsecase) AddUser(user *models.User) error{
	currentUsers, err := s.db.NumberOfUsers()

	if err != nil{
		return err
    }

    if currentUsers == 0{
        user.Role = "Admin"
    }
	
	return s.db.CreateUser(user)
}

func (s *UserUsecase) UpdateUser(id *string, updatedUser *models.User) (models.User, error) {
    // Call the repository's UpdateUser method
    updatedUserObj, err := s.db.UpdateUser(*id, updatedUser)
    if err != nil {
        return models.User{}, err
    }

    // Return the updated user object
    return updatedUserObj, nil
}


func (s *UserUsecase) GetAllUsers() ([]models.User, error) {
    users, err := s.db.GetUsers()
    if err != nil {
        return nil, err
    }
    return users, nil
}

