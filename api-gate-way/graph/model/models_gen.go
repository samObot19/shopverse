// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type FilterInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Mutation struct {
}

type Order struct {
	ID              string       `json:"id"`
	UserID          string       `json:"userID"`
	OrderStatus     string       `json:"orderStatus"`
	PaymentStatus   string       `json:"paymentStatus"`
	TotalAmount     float64      `json:"totalAmount"`
	ShippingAddress string       `json:"shippingAddress"`
	BillingAddress  string       `json:"billingAddress"`
	CreatedAt       string       `json:"createdAt"`
	UpdatedAt       string       `json:"updatedAt"`
	Items           []*OrderItem `json:"items"`
}

type OrderInput struct {
	UserID          string            `json:"userID"`
	Items           []*OrderItemInput `json:"items"`
	ShippingAddress string            `json:"shippingAddress"`
	BillingAddress  string            `json:"billingAddress"`
}

type OrderItem struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"productID"`
	ProductPrice float64 `json:"productPrice"`
	Quantity     int32   `json:"quantity"`
	TotalPrice   float64 `json:"totalPrice"`
}

type OrderItemInput struct {
	ProductID    string  `json:"productID"`
	ProductPrice float64 `json:"productPrice"`
	Quantity     int32   `json:"quantity"`
	TotalPrice   float64 `json:"totalPrice"`
}

type Product struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Stock       int32              `json:"stock"`
	Category    string             `json:"category"`
	Attributes  *ProductAttributes `json:"attributes"`
	Images      []string           `json:"images"`
	Ratings     float64            `json:"ratings"`
	CreatedAt   string             `json:"createdAt"`
}

type ProductAttributes struct {
	Color string `json:"color"`
	Size  string `json:"size"`
}

type ProductAttributesInput struct {
	Color string `json:"color"`
	Size  string `json:"size"`
}

type ProductInput struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Price       float64                 `json:"price"`
	Stock       int32                   `json:"stock"`
	Category    string                  `json:"category"`
	Attributes  *ProductAttributesInput `json:"attributes"`
	Images      []string                `json:"images"`
	Ratings     float64                 `json:"ratings"`
}

type Query struct {
}

type User struct {
	ID       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
}
