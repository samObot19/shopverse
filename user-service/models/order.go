package models

import "time"

type Order struct {
    ID             uint      `json:"id"`
    UserID         uint      `json:"user_id"`
    OrderStatus    string    `json:"order_status"`
    PaymentStatus  string    `json:"payment_status"`
    TotalAmount    float64   `json:"total_amount"`
    ShippingAddr   string    `json:"shipping_address"`
    BillingAddr    string    `json:"billing_address"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    Items          []OrderItem
}

type OrderItem struct {
    ID           uint    `json:"id"`
    OrderID      uint    `json:"order_id"`
    ProductID    uint    `json:"product_id"`
    ProductPrice float64 `json:"product_price"`
    Quantity     int     `json:"quantity"`
    TotalPrice   float64 `json:"total_price"`
}