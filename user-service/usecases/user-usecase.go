package usecase

import (
	"github.com/samObot19/shopverse/user-service/models"
	"github.com/samObot19/shopverse/user-service/repository"
	"github.com/samObot19/shopverse/user-service/events"
	"errors"
	"log"
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

func (s *UserUsecase) AddUser(user *models.User) error {
	
	if err := s.db.CreateUser(user); err != nil {
		return err
	}

	producer, err := events.NewUserEventProducer("localhost:9092", "user-events")
	if err != nil {
		log.Printf("Failed to initialize Kafka producer: %v", err)
		return err
	}
	defer producer.Close() 

	err = producer.PublishUserCreatedEvent(user)
	if err != nil {
		log.Printf("Failed to publish user created event: %v", err)
	}

	return nil
}

func (s *UserUsecase) UpdateUser(id *string, updatedUser *models.User) (models.User, error) {
    updatedUserObj, err := s.db.UpdateUser(*id, updatedUser)
    if err != nil {
        return models.User{}, err
    }
    return updatedUserObj, nil
}

func (s *UserUsecase) GetAllUsers() ([]models.User, error) {
    users, err := s.db.GetUsers()
    if err != nil {
        return nil, err
    }
    return users, nil
}

func (s *UserUsecase) GetUserByID(id string) (models.User, error) {
    user, err := s.db.GetUserByID(id)
    if err != nil {
        return models.User{}, err
    }
    return user, nil
}
