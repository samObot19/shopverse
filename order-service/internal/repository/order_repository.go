package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/samObot19/shopverse/order-service/internal/models"
)


type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) (uint, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	UpdatePaymentStatus(ctx context.Context, orderID string, status string) error
	DeleteOrder(ctx context.Context, orderID string) error
	GetAllOrders(ctx context.Context, userID string) ([]*models.Order, error)
}


type orderRepository struct {
	DB *sql.DB
}


func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{DB: db}
}


func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) (uint, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO orders (user_id, order_status, payment_status, total_amount, shipping_addr, billing_addr, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, order.UserID, order.OrderStatus, order.PaymentStatus, order.TotalAmount, order.ShippingAddr, order.BillingAddr, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	
	for _, item := range order.Items {
		itemQuery := `INSERT INTO order_items (order_id, product_id, product_price, quantity, total_price)
                      VALUES (?, ?, ?, ?, ?)`
		_, err := tx.Exec(itemQuery, orderID, item.ProductID, item.ProductPrice, item.Quantity, item.TotalPrice)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return uint(orderID), nil
}


func (r *orderRepository) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	query := `SELECT id, user_id, order_status, payment_status, total_amount, shipping_addr, billing_addr, created_at, updated_at
              FROM orders WHERE id = ?`
	row := r.DB.QueryRow(query, orderID)

	var order models.Order
	err := row.Scan(&order.ID, &order.UserID, &order.OrderStatus, &order.PaymentStatus, &order.TotalAmount, &order.ShippingAddr, &order.BillingAddr, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	
	itemQuery := `SELECT id, order_id, product_id, product_price, quantity, total_price FROM order_items WHERE order_id = ?`
	rows, err := r.DB.Query(itemQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductPrice, &item.Quantity, &item.TotalPrice)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}


func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := `UPDATE orders SET order_status = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, status, time.Now(), orderID)
	return err
}


func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, orderID string, status string) error {
	query := `UPDATE orders SET payment_status = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, status, time.Now(), orderID)
	return err
}


func (r *orderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	query := `DELETE FROM orders WHERE id = ?`
	_, err := r.DB.Exec(query, orderID)
	return err
}

// GetAllOrders retrieves all orders for a specific user
func (r *orderRepository) GetAllOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	query := `SELECT id, user_id, order_status, payment_status, total_amount, shipping_addr, billing_addr, created_at, updated_at
              FROM orders WHERE user_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderStatus, &order.PaymentStatus, &order.TotalAmount, &order.ShippingAddr, &order.BillingAddr, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Fetch order items for each order
		itemQuery := `SELECT id, order_id, product_id, product_price, quantity, total_price FROM order_items WHERE order_id = ?`
		itemRows, err := r.DB.QueryContext(ctx, itemQuery, order.ID)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		var items []models.OrderItem
		for itemRows.Next() {
			var item models.OrderItem
			err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductPrice, &item.Quantity, &item.TotalPrice)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		order.Items = items

		orders = append(orders, &order)
	}

	return orders, nil
}
