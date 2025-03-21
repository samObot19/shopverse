package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/samObot19/shopverse/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	CreateUser(data *models.User) error
	ReadUser(username string) (models.User, bool)
	ChangeRoleToAdmin(username string) error
	UpdateUser(username string, data *models.User) (models.User, error)
	NumberOfUsers() (int64, error)
	GetUsers() ([]models.User, error) // New method to get all users
}

type MongoUserRepo struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

func NewMongoStorage(mongoURL string)(*mongo.Client, error){
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
	
	if err != nil {
		return nil, err
    }
    
	return client, nil
}

func NewMongoUserRepo() *MongoUserRepo {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file:")
	}

	url := os.Getenv("url")

	client, err := NewMongoStorage(url)

	if err != nil {
		fmt.Println("here is the  what fuckkkkkk", err)
		fmt.Println(err)
		return nil
	}

	NewUsercollection := client.Database("Account").Collection("Users")

	return &MongoUserRepo{
		Collection: NewUsercollection,
		Client:     client,
	}
}

func (s *MongoUserRepo) CreateUser(data *models.User) error {
	_, err := s.Collection.InsertOne(context.TODO(), data)
	return err
}

func (s *MongoUserRepo) NumberOfUsers() (int64, error) {
	return s.Collection.CountDocuments(context.TODO(), bson.D{})
}

func (s *MongoUserRepo) ReadUser(username string) (models.User, bool) {
	var result models.User
	err := s.Collection.FindOne(context.TODO(), bson.M{"name": username}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, false
		}
		return models.User{}, false
	}
	return result, true
}

func (s *MongoUserRepo) UpdateUser(username string, data *models.User) (models.User, error) {
    filter := bson.M{"name": username}
    update := bson.M{"$set": bson.M{}}

    if data.Password != "" {
        update["$set"].(bson.M)["password"] = data.Password
    }
    if data.Role != "" {
        update["$set"].(bson.M)["role"] = data.Role
    }
    if data.Name != "" {
        update["$set"].(bson.M)["name"] = data.Name
    }
    if data.Email != "" {
        update["$set"].(bson.M)["email"] = data.Email
    }

    // Perform the update operation
    _, err := s.Collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return models.User{}, err
    }

    // Retrieve the updated user
    var updatedUser models.User
    err = s.Collection.FindOne(context.Background(), filter).Decode(&updatedUser)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return models.User{}, errors.New("user not found")
        }
        return models.User{}, err
    }

    return updatedUser, nil
}

func (s *MongoUserRepo) ChangeRoleToAdmin(username string) error {
	filter := bson.M{"name": username}
	update := bson.M{"$set": bson.M{"role": "Admin"}}

	result, err := s.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("the user does not exist")
	}

	return nil
}

// GetUsers retrieves all users from the collection
func (s *MongoUserRepo) GetUsers() ([]models.User, error) {
	var users []models.User

	// Find all documents in the collection
	cursor, err := s.Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into the users slice
	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for any errors during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

