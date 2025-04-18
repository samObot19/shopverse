package events

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/samObot19/shopverse/user-service/models"
)

type UserEventProducer struct {
    producer *kafka.Producer
    topic    string
}


func NewUserEventProducer(broker, topic string) (*UserEventProducer, error) {
    producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
    if err != nil {
        return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
    }

    return &UserEventProducer{
        producer: producer,
        topic:    topic,
    }, nil
}

// PublishUserCreatedEvent publishes a "user created" event to Kafka
func (p *UserEventProducer) PublishUserCreatedEvent(user *models.User) error {
    message, err := json.Marshal(user)
    if err != nil {
        return fmt.Errorf("failed to serialize user: %w", err)
    }

    // Produce the message to Kafka
    err = p.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
        Value:          message,
    }, nil)
    if err != nil {
        return fmt.Errorf("failed to produce Kafka message: %w", err)
    }

    log.Printf("User created event published for user ID: %s", user.ID.Hex())
    return nil
}

// PublishOrderEvent publishes an "order created" event with user and order details to Kafka
func (p *UserEventProducer) PublishOrderEvent(user *models.User, order *models.Order) error {
    payload := struct {
        User  *models.User  `json:"user"`
        Order *models.Order `json:"order"`
    }{
        User:  user,
        Order: order,
    }

   
    message, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to serialize order event: %w", err)
    }

    err = p.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
        Value:          message,
    }, nil)
    if err != nil {
        return fmt.Errorf("failed to produce Kafka message: %w", err)
    }

    log.Printf("Order event published for user ID: %s and order ID: %s", user.ID.Hex(), order.ID)
    return nil
}


func (p *UserEventProducer) Close() {
    p.producer.Close()
}