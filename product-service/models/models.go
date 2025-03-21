package models

import "time"

// Product represents the schema for a product
type Product struct {
    ID          string            `json:"_id" bson:"_id"`                     // Unique identifier
    Title       string            `json:"title" bson:"title"`                 // Product title
    Description string            `json:"description" bson:"description"`     // Product description
    Price       float64           `json:"price" bson:"price"`                 // Product price
    Stock       int               `json:"stock" bson:"stock"`                 // Stock quantity
    Category    string            `json:"category" bson:"category"`           // Product category
    Attributes  ProductAttributes `json:"attributes" bson:"attributes"`       // Additional attributes
    Images      []string          `json:"images" bson:"images"`               // List of image URLs
    Ratings     float64           `json:"ratings" bson:"ratings"`             // Product ratings
    CreatedAt   time.Time         `json:"created_at" bson:"created_at"`       // Creation timestamp
}

// ProductAttributes represents additional attributes for a product
type ProductAttributes struct {
    Color string   `json:"color" bson:"color"` // Product color
    Size  []string `json:"size" bson:"size"`   // Available sizes
}