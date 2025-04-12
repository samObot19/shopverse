package publish

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

const KafkaServer = "localhost:9092"

// PublishEvent publishes an event message to the specified Kafka topic

func PublishEvent(topic string, eventMessage interface{}) error {
    // Create a new Kafka producer
    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": KafkaServer,
    })
    if err != nil {
        return fmt.Errorf("failed to create Kafka producer: %v", err)
    }
    defer producer.Close()

    // Serialize the event message to JSON
    messageValue, err := json.Marshal(eventMessage)
    if err != nil {
        return fmt.Errorf("failed to serialize event message: %v", err)
    }

    // Produce the message
    err = producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          messageValue,
    }, nil)
    if err != nil {
        return fmt.Errorf("failed to produce message: %v", err)
    }

    // Wait for delivery event
    for e := range producer.Events() {
        switch ev := e.(type) {
        case *kafka.Message:
            if ev.TopicPartition.Error != nil {
                return fmt.Errorf("delivery failed: %v", ev.TopicPartition.Error)
            }
            fmt.Printf("✅ Event published to topic %s: %s\n", topic, string(messageValue))
            return nil
        case kafka.Error:
            // log error and return, or optionally retry
            return fmt.Errorf("kafka error: %v", ev)
        default:
            // other events can be ignored or logged
            log.Printf("⚠️ Ignored event type: %T\n", ev)
        }
    }

    return fmt.Errorf("no delivery event received")
}
