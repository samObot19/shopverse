package orderclient

import (
	"context"
	"log"
	"time"

	pb "github.com/samObot19/shopverse/api-gate-way/order-client/proto/pb"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client pb.OrderServiceClient
}

// NewOrderClient creates a new OrderClient
func NewOrderClient(conn *grpc.ClientConn) *OrderClient {
	return &OrderClient{
		client: pb.NewOrderServiceClient(conn),
	}
}

// CreateOrder calls the CreateOrder gRPC method
func (oc *OrderClient) CreateOrder(ctx context.Context, userID uint32, items []*pb.OrderItem, shippingAddress, billingAddress string) (*pb.CreateOrderResponse, error) {
	req := &pb.CreateOrderRequest{
		UserId:          userID,
		Items:           items,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
	}
	resp, err := oc.client.CreateOrder(ctx, req)
	if (err != nil) {
		log.Printf("Error creating order: %v", err)
		return nil, err
	}
	log.Println("Order created successfully")
	return resp, nil
}

// GetOrderByID calls the GetOrderByID gRPC method
func (oc *OrderClient) GetOrderByID(ctx context.Context, orderID uint32) (*pb.Order, error) {
	req := &pb.GetOrderByIDRequest{
		OrderId: orderID,
	}
	resp, err := oc.client.GetOrderByID(ctx, req)
	if err != nil {
		log.Printf("Error fetching order by ID: %v", err)
		return nil, err
	}
	return resp.Order, nil
}

// UpdateOrderStatus calls the UpdateOrderStatus gRPC method
func (oc *OrderClient) UpdateOrderStatus(ctx context.Context, orderID uint32, status string) (*pb.UpdateOrderStatusResponse, error) {
	req := &pb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  status,
	}
	resp, err := oc.client.UpdateOrderStatus(ctx, req)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		return nil, err
	}
	log.Println("Order status updated successfully")
	return resp, nil
}

// UpdatePaymentStatus calls the UpdatePaymentStatus gRPC method
func (oc *OrderClient) UpdatePaymentStatus(ctx context.Context, orderID uint32, paymentStatus string) (*pb.UpdatePaymentStatusResponse, error) {
	req := &pb.UpdatePaymentStatusRequest{
		OrderId:       orderID,
		PaymentStatus: paymentStatus,
	}
	resp, err := oc.client.UpdatePaymentStatus(ctx, req)
	if err != nil {
		log.Printf("Error updating payment status: %v", err)
		return nil, err
	}
	log.Println("Payment status updated successfully")
	return resp, nil
}

// DeleteOrder calls the DeleteOrder gRPC method
func (oc *OrderClient) DeleteOrder(ctx context.Context, orderID uint32) (*pb.DeleteOrderResponse, error) {
	req := &pb.DeleteOrderRequest{
		OrderId: orderID,
	}
	resp, err := oc.client.DeleteOrder(ctx, req)
	if err != nil {
		log.Printf("Error deleting order: %v", err)
		return nil, err
	}
	log.Println("Order deleted successfully")
	return resp, nil
}

// GetAllOrders calls the GetAllOrders gRPC method
func (oc *OrderClient) GetAllOrders(ctx context.Context, userID uint32) ([]*pb.Order, error) {
	req := &pb.GetAllOrdersRequest{
		UserId: userID,
	}
	resp, err := oc.client.GetAllOrders(ctx, req)
	if err != nil {
		log.Printf("Error fetching all orders: %v", err)
		return nil, err
	}
	return resp.Orders, nil
}

// ConnectToOrderService establishes a connection to the order service
func ConnectToOrderService(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Printf("Failed to connect to order service: %v", err)
		return nil, err
	}
	return conn, nil
}
