package service

import (
    "context"
    "time"

    pb "github.com/samObot19/shopverse/product-service/proto/pb"
    "github.com/samObot19/shopverse/product-service/usecases"
    "github.com/samObot19/shopverse/product-service/models"
)

type ProductServiceServer struct {
    pb.UnimplementedProductServiceServer
    useCase *usecases.ProductUseCase
}

// NewProductServiceServer creates a new ProductServiceServer
func NewProductServiceServer(useCase *usecases.ProductUseCase) *ProductServiceServer {
    return &ProductServiceServer{useCase: useCase}
}

// Helper functions for time conversion
func parseTime(timeStr string) (time.Time, error) {
    return time.Parse(time.RFC3339, timeStr)
}

func formatTime(t time.Time) string {
    return t.Format(time.RFC3339)
}

// CreateProduct handles the gRPC request to create a new product
func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
    createdAt, err := parseTime(req.Product.CreatedAt)
    if err != nil {
        return nil, err
    }

    product := &models.Product{
        ID:          req.Product.Id,
        Title:       req.Product.Title,
        Description: req.Product.Description,
        Price:       req.Product.Price,
        Stock:       int(req.Product.Stock),
        Category:    req.Product.Category,
        Attributes: models.ProductAttributes{
            Color: req.Product.Attributes.Color,
            Size:  req.Product.Attributes.Size,
        },
        Images:    req.Product.Images,
        Ratings:   req.Product.Ratings,
        CreatedAt: createdAt,
    }

    err = s.useCase.CreateProduct(ctx, product)
    if err != nil {
        return nil, err
    }

    return &pb.CreateProductResponse{Message: "Product created successfully"}, nil
}

// GetProductByID handles the gRPC request to retrieve a product by its ID
func (s *ProductServiceServer) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
    product, err := s.useCase.GetProductByID(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    return &pb.GetProductByIDResponse{
        Product: &pb.Product{
            Id:          product.ID,
            Title:       product.Title,
            Description: product.Description,
            Price:       product.Price,
            Stock:       int32(product.Stock),
            Category:    product.Category,
            Attributes: &pb.Attributes{
                Color: product.Attributes.Color,
                Size:  product.Attributes.Size,
            },
            Images:    product.Images,
            Ratings:   product.Ratings,
            CreatedAt: formatTime(product.CreatedAt),
        },
    }, nil
}

// GetAllProducts handles the gRPC request to retrieve all products
func (s *ProductServiceServer) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
    filters := make(map[string]interface{})
    for key, value := range req.Filters {
        filters[key] = value
    }

    products, err := s.useCase.GetAllProducts(ctx, filters)
    if err != nil {
        return nil, err
    }

    var pbProducts []*pb.Product
    for _, product := range products {
        pbProducts = append(pbProducts, &pb.Product{
            Id:          product.ID,
            Title:       product.Title,
            Description: product.Description,
            Price:       product.Price,
            Stock:       int32(product.Stock),
            Category:    product.Category,
            Attributes: &pb.Attributes{
                Color: product.Attributes.Color,
                Size:  product.Attributes.Size,
            },
            Images:    product.Images,
            Ratings:   product.Ratings,
            CreatedAt: formatTime(product.CreatedAt),
        })
    }

    return &pb.GetAllProductsResponse{Products: pbProducts}, nil
}

// UpdateProduct handles the gRPC request to update an existing product
func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
    product := &models.Product{
        ID:          req.Product.Id,
        Title:       req.Product.Title,
        Description: req.Product.Description,
        Price:       req.Product.Price,
        Stock:       int(req.Product.Stock),
        Category:    req.Product.Category,
        Attributes: models.ProductAttributes{
            Color: req.Product.Attributes.Color,
            Size:  req.Product.Attributes.Size,
        },
        Images:    req.Product.Images,
        Ratings:   req.Product.Ratings,
        CreatedAt: func() time.Time {
            parsedTime, _ := parseTime(req.Product.CreatedAt)
            return parsedTime
        }(),
    }

    err := s.useCase.UpdateProduct(ctx, req.Id, product)
    if err != nil {
        return nil, err
    }

    return &pb.UpdateProductResponse{Message: "Product updated successfully"}, nil
}

// DeleteProduct handles the gRPC request to delete a product by its ID
func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
    err := s.useCase.DeleteProduct(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    return &pb.DeleteProductResponse{Message: "Product deleted successfully"}, nil
}

// UpdateStock handles the gRPC request to update the stock quantity of a product
func (s *ProductServiceServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
    err := s.useCase.UpdateStock(ctx, req.Id, int(req.Quantity))
    if err != nil {
        return nil, err
    }

    return &pb.UpdateStockResponse{Message: "Stock updated successfully"}, nil
}

// GetProductsByCategory handles the gRPC request to retrieve products by category
func (s *ProductServiceServer) GetProductsByCategory(ctx context.Context, req *pb.GetProductsByCategoryRequest) (*pb.GetProductsByCategoryResponse, error) {
    products, err := s.useCase.GetProductsByCategory(ctx, req.Category)
    if err != nil {
        return nil, err
    }

    var pbProducts []*pb.Product
    for _, product := range products {
        pbProducts = append(pbProducts, &pb.Product{
            Id:          product.ID,
            Title:       product.Title,
            Description: product.Description,
            Price:       product.Price,
            Stock:       int32(product.Stock),
            Category:    product.Category,
            Attributes: &pb.Attributes{
                Color: product.Attributes.Color,
                Size:  product.Attributes.Size,
            },
            Images:    product.Images,
            Ratings:   product.Ratings,
            CreatedAt: formatTime(product.CreatedAt),
        })
    }

    return &pb.GetProductsByCategoryResponse{Products: pbProducts}, nil
}

// SearchProducts handles the gRPC request to search for products
func (s *ProductServiceServer) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
    products, err := s.useCase.SearchProducts(ctx, req.Query)
    if err != nil {
        return nil, err
    }

    var pbProducts []*pb.Product
    for _, product := range products {
        pbProducts = append(pbProducts, &pb.Product{
            Id:          product.ID,
            Title:       product.Title,
            Description: product.Description,
            Price:       product.Price,
            Stock:       int32(product.Stock),
            Category:    product.Category,
            Attributes: &pb.Attributes{
                Color: product.Attributes.Color,
                Size:  product.Attributes.Size,
            },
            Images:    product.Images,
            Ratings:   product.Ratings,
            CreatedAt: formatTime(product.CreatedAt),
        })
    }

    return &pb.SearchProductsResponse{Products: pbProducts}, nil
}