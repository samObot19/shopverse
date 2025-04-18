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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	CreateUser(data *models.User) error
	ReadUser(username string) (models.User, bool)
	ChangeRoleToAdmin(username string) error
	NumberOfUsers() (int64, error)
	UpdateUser(username string, data *models.User) (models.User, error)
	GetUsers() ([]models.User, error) 
	GetUserByID(id string) (models.User, error)
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


func (s *MongoUserRepo) ReadUser(email string) (models.User, bool) {
	var result models.User
	err := s.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)

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
    if data.ProfilePicture != "" {
        update["$set"].(bson.M)["profile_picture"] = data.ProfilePicture
    }

  
    if len(update["$set"].(bson.M)) == 0 {
        return models.User{}, errors.New("no fields to update")
    }

    _, err := s.Collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return models.User{}, err
    }

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


func (s *MongoUserRepo) GetUsers() ([]models.User, error) {
	var users []models.User

	
	cursor, err := s.Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}


func (s *MongoUserRepo) GetUserByID(id string) (models.User, error) {
    var user models.User

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return models.User{}, errors.New("invalid ObjectId format")
    }

    err = s.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return models.User{}, errors.New("user not found")
        }
        return models.User{}, err
    }

    return user, nil
}

