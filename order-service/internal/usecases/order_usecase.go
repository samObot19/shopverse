package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"github.com/samObot19/shopverse/order-service/internal/events/publish"
	"github.com/samObot19/shopverse/order-service/internal/models"
	"github.com/samObot19/shopverse/order-service/internal/repository"
)

// OrderUsecase defines the interface for order-related business logic.
type OrderUsecase interface {
	CreateOrder(ctx context.Context, order *models.Order) (uint, error)
	GetOrderByID(ctx context.Context, orderID uint) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	UpdatePaymentStatus(ctx context.Context, orderID uint, status string) error
	DeleteOrder(ctx context.Context, orderID uint) error
	GetAllOrders(ctx context.Context, userID uint) ([]*models.Order, error)
}

// orderUsecase is a concrete implementation of the OrderUsecase interface.
type orderUsecase struct {
	repo repository.OrderRepository
}

// NewOrderUsecase creates a new instance of orderUsecase.
func NewOrderUsecase(repo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{repo: repo}
}

// CreateOrder handles the business logic for creating an order.
func (u *orderUsecase) CreateOrder(ctx context.Context, order *models.Order) (uint, error) {
	// Validate order data
	if len(order.Items) == 0 {
		return 0, errors.New("order must contain at least one item")
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range order.Items {
		item.TotalPrice = item.ProductPrice * float64(item.Quantity)
		totalAmount += item.TotalPrice
	}
	order.TotalAmount = totalAmount

	// Set default order and payment statuses
	order.OrderStatus = "Pending"
	order.PaymentStatus = "Unpaid"

	// Call repository to create the order
	orderID, err := u.repo.CreateOrder(ctx, order)
	fmt.Println("order created")
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return 0, err
	}

	// Publish the order-created event to the orderCreated topic
	err = publish.PublishEvent("orderCreated", *order)
	if err != nil {
		log.Printf("Failed to publish order-created event: %v", err)
		return 0, err
	}

	return orderID, nil
}

// GetOrderByID handles the business logic for retrieving an order by its ID.
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

// UpdateOrderStatus handles the business logic for updating the order status.
func (u *orderUsecase) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	// Validate status
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

	// Call repository to update the status
	err := u.repo.UpdateOrderStatus(ctx, fmt.Sprint("%d", orderID), status)
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		return err
	}
	return nil
}

// UpdatePaymentStatus handles the business logic for updating the payment status.
func (u *orderUsecase) UpdatePaymentStatus(ctx context.Context, orderID uint, status string) error {
	// Validate status
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

	// Call repository to update the payment status
	err := u.repo.UpdatePaymentStatus(ctx, fmt.Sprintf("%d", orderID), status)
	if err != nil {
		log.Printf("Failed to update payment status: %v", err)
		return err
	}
	return nil
}

// DeleteOrder handles the business logic for deleting an order.
func (u *orderUsecase) DeleteOrder(ctx context.Context, orderID uint) error {
	err := u.repo.DeleteOrder(ctx, fmt.Sprintf("%d",orderID))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		return err
	}
	return nil
}

// GetAllOrders handles the business logic for retrieving all orders for a user.
func (u *orderUsecase) GetAllOrders(ctx context.Context, userID uint) ([]*models.Order, error) {
	orders, err := u.repo.GetAllOrders(ctx, userID)
	if err != nil {
		log.Printf("Failed to retrieve orders: %v", err)
		return nil, err
	}
	return orders, nil
}
