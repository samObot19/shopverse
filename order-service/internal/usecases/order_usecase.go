package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"github.com/samObot19/shopverse/order-service/internal/events/publish"
	"github.com/samObot19/shopverse/order-service/internal/models"
	"github.com/samObot19/shopverse/order-service/internal/repository"
	"github.com/samObot19/shopverse/order-service/clients/product-client/proto/pb"
)


type OrderUsecase interface {
	CreateOrder(ctx context.Context, order *models.Order) (uint, error)
	GetOrderByID(ctx context.Context, orderID uint) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	UpdatePaymentStatus(ctx context.Context, orderID uint, status string) error
	DeleteOrder(ctx context.Context, orderID uint) error
	GetAllOrders(ctx context.Context, userID string) ([]*models.Order, error)
}


type orderUsecase struct {
	repo          repository.OrderRepository
	productClient pb.ProductServiceClient
}


func NewOrderUsecase(repo repository.OrderRepository, productClient pb.ProductServiceClient) OrderUsecase {
	return &orderUsecase{
		repo:          repo,
		productClient: productClient,
	}
}


func (u *orderUsecase) CreateOrder(ctx context.Context, order *models.Order) (uint, error) {
	if len(order.Items) == 0 {
		return 0, errors.New("order must contain at least one item")
	}

	for _, item := range order.Items {
		productResponse, err := u.productClient.GetProductByID(ctx, &pb.GetProductByIDRequest{
			Id: item.ProductID,
		})
		if err != nil {
			log.Printf("Failed to fetch product details for product ID %s: %v", item.ProductID, err)
			return 0, fmt.Errorf("failed to fetch product details for product ID %s", item.ProductID)
		}

		product := productResponse.Product
		if product.Stock < int32(item.Quantity) {
			log.Printf("Insufficient stock for product ID %s: available %d, required %d", item.ProductID, product.Stock, item.Quantity)
			return 0, fmt.Errorf("insufficient stock for product ID %s", item.ProductID)
		}
	}

	var totalAmount float64
	for _, item := range order.Items {
		item.TotalPrice = item.ProductPrice * float64(item.Quantity)
		totalAmount += item.TotalPrice
	}
	order.TotalAmount = totalAmount

	order.OrderStatus = "Pending"
	order.PaymentStatus = "Unpaid"

	orderID, err := u.repo.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return 0, err
	}

	err = publish.PublishEvent("orderCreated", *order)
	if err != nil {
		log.Printf("Failed to publish order-created event: %v", err)
		return 0, err
	}

	return orderID, nil
}


func (u *orderUsecase) GetOrderByID(ctx context.Context, orderID uint) (*models.Order, error) {
	order, err := u.repo.GetOrderByID(ctx, fmt.Sprintf("%d", orderID))
	if err != nil {
		log.Printf("Failed to retrieve order: %v", err)
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}


func (u *orderUsecase) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	validStatuses := []string{"Pending", "Processing", "Shipped", "Delivered", "Cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid order status")
	}

	err := u.repo.UpdateOrderStatus(ctx, fmt.Sprintf("%d", orderID), status)
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		return err
	}
	return nil
}


func (u *orderUsecase) UpdatePaymentStatus(ctx context.Context, orderID uint, status string) error {
	validStatuses := []string{"Unpaid", "Paid", "Refunded"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid payment status")
	}

	err := u.repo.UpdatePaymentStatus(ctx, fmt.Sprintf("%d", orderID), status)
	if err != nil {
		log.Printf("Failed to update payment status: %v", err)
		return err
	}
	return nil
}

func (u *orderUsecase) DeleteOrder(ctx context.Context, orderID uint) error {
	err := u.repo.DeleteOrder(ctx, fmt.Sprintf("%d",orderID))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		return err
	}
	return nil
}


func (u *orderUsecase) GetAllOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	orders, err := u.repo.GetAllOrders(ctx, userID)
	if err != nil {
		log.Printf("Failed to retrieve orders: %v", err)
		return nil, err
	}
	return orders, nil
}
