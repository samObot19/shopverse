package services

import (
	"context"
	"log"

	"github.com/samObot19/shopverse/order-service/internal/models"
	"github.com/samObot19/shopverse/order-service/internal/usecases"
	"github.com/samObot19/shopverse/order-service/proto/pb"
)


type OrderServiceServer struct {
	usecase usecases.OrderUsecase
	pb.UnimplementedOrderServiceServer
}

func NewOrderServiceServer(usecase usecases.OrderUsecase) *OrderServiceServer {
	return &OrderServiceServer{usecase: usecase}
}


func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &models.Order{
		UserID:       req.UserId, 
		ShippingAddr: req.ShippingAddress,
		BillingAddr:  req.BillingAddress,
		Items:        convertProtoOrderItemsToModel(req.Items),
	}

	orderID, err := s.usecase.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, err
	}

	return &pb.CreateOrderResponse{OrderId: uint32(orderID)}, nil
}


func (s *OrderServiceServer) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	order, err := s.usecase.GetOrderByID(ctx, uint(req.OrderId))
	if err != nil {
		log.Printf("Failed to retrieve order: %v", err)
		return nil, err
	}

	return &pb.GetOrderByIDResponse{Order: convertModelOrderToProto(order)}, nil
}


func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	err := s.usecase.UpdateOrderStatus(ctx, uint(req.OrderId), req.Status)
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		return nil, err
	}

	return &pb.UpdateOrderStatusResponse{Message: "Order status updated successfully"}, nil
}


func (s *OrderServiceServer) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error) {
	err := s.usecase.UpdatePaymentStatus(ctx, uint(req.OrderId), req.PaymentStatus)
	if err != nil {
		log.Printf("Failed to update payment status: %v", err)
		return nil, err
	}

	return &pb.UpdatePaymentStatusResponse{Message: "Payment status updated successfully"}, nil
}


func (s *OrderServiceServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := s.usecase.DeleteOrder(ctx, uint(req.OrderId))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		return nil, err
	}

	return &pb.DeleteOrderResponse{Message: "Order deleted successfully"}, nil
}


func (s *OrderServiceServer) GetAllOrders(ctx context.Context, req *pb.GetAllOrdersRequest) (*pb.GetAllOrdersResponse, error) {
	orders, err := s.usecase.GetAllOrders(ctx, req.UserId) 
	if err != nil {
		log.Printf("Failed to retrieve orders: %v", err)
		return nil, err
	}

	return &pb.GetAllOrdersResponse{Orders: convertModelOrdersToProto(orders)}, nil
}


func convertProtoOrderItemsToModel(items []*pb.OrderItem) []models.OrderItem {
	var modelItems []models.OrderItem
	for _, item := range items {
		modelItems = append(modelItems, models.OrderItem{
			ProductID:    item.ProductId, 
			ProductPrice: float64(item.ProductPrice),
			Quantity:     int(item.Quantity),
			TotalPrice:   float64(item.TotalPrice),
		})
	}
	return modelItems
}

func convertModelOrderToProto(order *models.Order) *pb.Order {
	return &pb.Order{
		Id:              uint32(order.ID),
		UserId:          order.UserID, 
		OrderStatus:     order.OrderStatus,
		PaymentStatus:   order.PaymentStatus,
		TotalAmount:     float32(order.TotalAmount),
		ShippingAddress: order.ShippingAddr,
		BillingAddress:  order.BillingAddr,
		CreatedAt:       order.CreatedAt.String(),
		UpdatedAt:       order.UpdatedAt.String(),
		Items:           convertModelOrderItemsToProto(order.Items),
	}
}

func convertModelOrderItemsToProto(items []models.OrderItem) []*pb.OrderItem {
	var protoItems []*pb.OrderItem
	for _, item := range items {
		protoItems = append(protoItems, &pb.OrderItem{
			Id:           uint32(item.ID),
			ProductId:    item.ProductID, 
			ProductPrice: float32(item.ProductPrice),
			Quantity:     uint32(item.Quantity),
			TotalPrice:   float32(item.TotalPrice),
		})
	}
	return protoItems
}

func convertModelOrdersToProto(orders []*models.Order) []*pb.Order {
	var protoOrders []*pb.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, convertModelOrderToProto(order))
	}
	return protoOrders
}
