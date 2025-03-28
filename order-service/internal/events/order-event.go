package order_event

import (
    "encoding/json"
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/samObot19/shopverse/order-service/internal/models"
)

const KafkaServer = "localhost:9092"

var KafkaTopic = "orders-v1-topic"

func PublishOrderEvent(order models.Order) error {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": KafkaServer,
    })
    if err != nil {
        return fmt.Errorf("failed to create Kafka producer: %v", err)
    }
    defer p.Close()

    value, err := json.Marshal(order)
    if err != nil {
        return fmt.Errorf("failed to serialize order: %v", err)
    }

    err = p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &KafkaTopic, Partition: kafka.PartitionAny},
        Value:          value,
    }, nil)
    if err != nil {
        return fmt.Errorf("failed to produce message: %v", err)
    }

    e := <-p.Events()
    m := e.(*kafka.Message)
    if m.TopicPartition.Error != nil {
        return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
    }

    fmt.Printf("Order event published to topic %s: %s\n", KafkaTopic, string(value))
    return nil
}

