package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // MongoDB ObjectId
    GoogleID       string             `json:"google_id" bson:"google_id"`        // Google unique user ID
    Name           string             `json:"name" bson:"name"`                 // User's full name
    Email          string             `json:"email" bson:"email"`               // User's email address
	Password 	   string             `json:"password" bson:"password"`
    ProfilePicture string             `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"` // URL to the user's profile picture
    Role           string             `json:"role" bson:"role"`                 // User role (e.g., "customer", "admin")
}
