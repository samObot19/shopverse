package models

type OrderMessage struct {
    OrderID     string `json:"order_id"`
    UserID      string `json:"user_id"`
    UserEmail   string `json:"user_email"`
    ProductID   string `json:"product_id"`
    Quantity    int    `json:"quantity"`
    TotalAmount float64 `json:"total_amount"`
    CreatedAt   string `json:"created_at"`
}