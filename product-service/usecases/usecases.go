package usecases

import (
    "context"
    "errors"
    "github.com/samObot19/shopverse/product-service/models"
    "github.com/samObot19/shopverse/product-service/repository"
)

type ProductUseCase struct {
    repo repository.ProductRepository
}

// NewProductUseCase creates a new instance of ProductUseCase
func NewProductUseCase(repo repository.ProductRepository) *ProductUseCase {
    return &ProductUseCase{repo: repo}
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *models.Product) error {
    if product.Title == "" || product.Price <= 0 || product.Stock < 0 {
        return errors.New("invalid product data")
    }
    return uc.repo.CreateProduct(ctx, product)
}

// GetProductByID retrieves a product by its ID
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
    if id == "" {
        return nil, errors.New("product ID cannot be empty")
    }
    return uc.repo.GetProductByID(ctx, id)
}

// GetAllProducts retrieves all products with optional filters
func (uc *ProductUseCase) GetAllProducts(ctx context.Context, filters map[string]interface{}) ([]*models.Product, error) {
    return uc.repo.GetAllProducts(ctx, filters)
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id string, updatedProduct *models.Product) error {
    if id == "" {
        return errors.New("product ID cannot be empty")
    }
    return uc.repo.UpdateProduct(ctx, id, updatedProduct)
}

// DeleteProduct deletes a product by its ID
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("product ID cannot be empty")
    }
    return uc.repo.DeleteProduct(ctx, id)
}

// UpdateStock updates the stock quantity of a product
func (uc *ProductUseCase) UpdateStock(ctx context.Context, id string, quantity int) error {
    if id == "" {
        return errors.New("product ID cannot be empty")
    }
    if quantity == 0 {
        return errors.New("quantity must not be zero")
    }
    return uc.repo.UpdateStock(ctx, id, quantity)
}

// GetProductsByCategory retrieves products by category
func (uc *ProductUseCase) GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error) {
    if category == "" {
        return nil, errors.New("category cannot be empty")
    }
    return uc.repo.GetProductsByCategory(ctx, category)
}

// SearchProducts searches for products based on a query string
func (uc *ProductUseCase) SearchProducts(ctx context.Context, query string) ([]*models.Product, error) {
    if query == "" {
        return nil, errors.New("search query cannot be empty")
    }
    return uc.repo.SearchProducts(ctx, query)
}