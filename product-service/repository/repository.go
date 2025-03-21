package repository

import (
    "context"
    "github.com/samObot19/shopverse/product-service/models"
)

// ProductRepository defines the interface for product-related database operations
type ProductRepository interface {
    CreateProduct(ctx context.Context, product *models.Product) error
    GetProductByID(ctx context.Context, id string) (*models.Product, error)
    GetAllProducts(ctx context.Context, filters map[string]interface{}) ([]*models.Product, error)
    UpdateProduct(ctx context.Context, id string, updatedProduct *models.Product) error
    DeleteProduct(ctx context.Context, id string) error
    UpdateStock(ctx context.Context, id string, quantity int) error
    GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error)
    SearchProducts(ctx context.Context, query string) ([]*models.Product, error)
}