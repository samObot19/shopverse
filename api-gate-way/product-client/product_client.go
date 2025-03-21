package productclient

import (
	"context"
	"log"
	"strings"
	"time"

	pb "github.com/samObot19/shopverse/api-gate-way/product-client/proto/pb"
	"github.com/samObot19/shopverse/api-gate-way/graph/model"
	"google.golang.org/grpc"
)

type ProductClient struct {
	client pb.ProductServiceClient
}

// NewProductClient creates a new ProductClient
func NewProductClient(conn *grpc.ClientConn) *ProductClient {
	return &ProductClient{
		client: pb.NewProductServiceClient(conn),
	}
}

// CreateProduct calls the CreateProduct gRPC method
func (pc *ProductClient) CreateProduct(ctx context.Context, product *pb.Product) error {
	_, err := pc.client.CreateProduct(ctx, &pb.CreateProductRequest{Product: product})
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return err
	}
	log.Println("Product created successfully")
	return nil
}

// GetProductByID calls the GetProductByID gRPC method
func (pc *ProductClient) GetProductByID(ctx context.Context, id string) (*pb.Product, error) {
	resp, err := pc.client.GetProductByID(ctx, &pb.GetProductByIDRequest{Id: id})
	if err != nil {
		log.Printf("Error fetching product by ID: %v", err)
		return nil, err
	}
	return resp.Product, nil
}

// GetAllProducts calls the GetAllProducts gRPC method
func (pc *ProductClient) GetAllProducts(ctx context.Context, filters map[string]string) ([]*pb.Product, error) {
	resp, err := pc.client.GetAllProducts(ctx, &pb.GetAllProductsRequest{Filters: filters})
	if err != nil {
		log.Printf("Error fetching all products: %v", err)
		return nil, err
	}
	return resp.Products, nil
}

// UpdateProduct calls the UpdateProduct gRPC method
func (pc *ProductClient) UpdateProduct(ctx context.Context, id string, product *pb.Product) error {
	_, err := pc.client.UpdateProduct(ctx, &pb.UpdateProductRequest{Id: id, Product: product})
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return err
	}
	log.Println("Product updated successfully")
	return nil
}

// DeleteProduct calls the DeleteProduct gRPC method
func (pc *ProductClient) DeleteProduct(ctx context.Context, id string) error {
	_, err := pc.client.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return err
	}
	log.Println("Product deleted successfully")
	return nil
}

// UpdateStock calls the UpdateStock gRPC method
func (pc *ProductClient) UpdateStock(ctx context.Context, id string, quantity int32) error {
	_, err := pc.client.UpdateStock(ctx, &pb.UpdateStockRequest{Id: id, Quantity: quantity})
	if err != nil {
		log.Printf("Error updating stock: %v", err)
		return err
	}
	log.Println("Stock updated successfully")
	return nil
}

// GetProductsByCategory calls the GetProductsByCategory gRPC method
func (pc *ProductClient) GetProductsByCategory(ctx context.Context, category string) ([]*pb.Product, error) {
	resp, err := pc.client.GetProductsByCategory(ctx, &pb.GetProductsByCategoryRequest{Category: category})
	if err != nil {
		log.Printf("Error fetching products by category: %v", err)
		return nil, err
	}
	return resp.Products, nil
}

// SearchProducts calls the SearchProducts gRPC method
func (pc *ProductClient) SearchProducts(ctx context.Context, query string) ([]*pb.Product, error) {
	resp, err := pc.client.SearchProducts(ctx, &pb.SearchProductsRequest{Query: query})
	if err != nil {
		log.Printf("Error searching products: %v", err)
		return nil, err
	}
	return resp.Products, nil
}

// ConnectToProductService establishes a connection to the product service
func ConnectToProductService(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Printf("Failed to connect to product service: %v", err)
		return nil, err
	}
	return conn, nil
}

type Product struct {
	ID          string
	Title       string
	Description string
	Price       float64
	Stock       int32
	Category    string
	Attributes  *ProductAttributes
	Images      []string
	Ratings     float64
	CreatedAt   string
}

type ProductAttributes struct {
	Color string
	Size  string
}

// FromProtoProduct converts a pb.Product to a model.Product
func FromProtoProduct(protoProduct *pb.Product) *model.Product {
	if protoProduct == nil {
		return nil
	}

	return &model.Product{
		ID:          protoProduct.Id,
		Title:       protoProduct.Title,
		Description: protoProduct.Description,
		Price:       protoProduct.Price,
		Stock:       int32(protoProduct.Stock), // Fixed type mismatch
		Category:    protoProduct.Category,
		Attributes: &model.ProductAttributes{
			Color: protoProduct.Attributes.Color,
			Size:  strings.Join(protoProduct.Attributes.Size, ","),
		},
		Images:    protoProduct.Images,
		Ratings:   protoProduct.Ratings,
		CreatedAt: protoProduct.CreatedAt,
	}
}

// ToProtoProduct converts a model.Product to a pb.Product
func ToProtoProduct(product *model.Product) *pb.Product {
	if product == nil {
		return nil
	}

	return &pb.Product{
		Id:          product.ID, // Ensure ID is passed
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
		Category:    product.Category,
		Attributes: &pb.Attributes{
			Color: product.Attributes.Color,
			Size:  strings.Split(product.Attributes.Size, ","),
		},
		Images:    product.Images,
		Ratings:   product.Ratings,
		CreatedAt: product.CreatedAt,
	}
}
