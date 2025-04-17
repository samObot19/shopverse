package subscribe

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/samObot19/shopverse/order-service/internal/events/publish"
    "github.com/samObot19/shopverse/order-service/internal/usecases"
)

const (
    KafkaServer       = "localhost:9092"
    StockEventTopic   = "stockEvent"
    OrderEventTopic   = "orderEvent"
    GroupID           = "order-service-group"
)

// SubscribeAndProcessStockEvent subscribes to the stockEvent topic and publishes to the orderEvent topic
func SubscribeAndProcessStockEvent(orderUsecase usecases.OrderUsecase) {
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": KafkaServer,
        "group.id":          GroupID,
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        panic(fmt.Sprintf("Failed to create Kafka consumer: %v", err))
    }
    defer consumer.Close()

    // Subscribe to the stockEvent topic
    err = consumer.SubscribeTopics([]string{StockEventTopic}, nil)
    if err != nil {
        panic(fmt.Sprintf("Failed to subscribe to topic: %v", err))
    }

    fmt.Println("Listening to stockEvent topic...")
    for {
        // Read messages from the stockEvent topic
        msg, err := consumer.ReadMessage(-1)
        if err == nil {
            fmt.Printf("Received stock event: %s\n", string(msg.Value))

            // Deserialize the stock event message
            var stockEvent struct {
                OrderID   uint `json:"order_id"`
                StockAvailable bool `json:"stock_available"`
            }
            err := json.Unmarshal(msg.Value, &stockEvent)
            if err != nil {
                fmt.Printf("Failed to deserialize stock event: %v\n", err)
                continue
            }

            // Update order status based on stock availability
            var newStatus string
            if stockEvent.StockAvailable {
                newStatus = "Accepted"
            } else {
                newStatus = "Failed"
            }

            err = orderUsecase.UpdateOrderStatus(context.Background(), stockEvent.OrderID, newStatus)
            if err != nil {
                fmt.Printf("Failed to update order status: %v\n", err)
                continue
            }

            // Publish the updated order to the orderEvent topic
            order, err := orderUsecase.GetOrderByID(context.Background(), stockEvent.OrderID)
            if err != nil {
                fmt.Printf("Failed to retrieve updated order: %v\n", err)
                continue
            }

            err = publish.PublishEvent(OrderEventTopic, order)
            if err != nil {
                fmt.Printf("Failed to publish order event: %v\n", err)
            } else {
                fmt.Println("Successfully published updated order to orderEvent topic.")
            }
        } else {
            fmt.Printf("Consumer error: %v\n", err)
        }
    }
}



