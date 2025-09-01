package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

// RawEvent represents the event published by the API service
type RawEvent struct {
	UserID    string `json:"user_id"`
	EventType string `json:"event_type"`
	Timestamp int64  `json:"timestamp"`
	Meta      string `json:"meta"`
}

// ProcessedEvent represents the enriched version of RawEvent
type ProcessedEvent struct {
	UserID       string `json:"user_id"`
	EventType    string `json:"event_type"`
	Timestamp    int64  `json:"timestamp"`
	Device       string `json:"device"`
	GeoLocation  string `json:"geo_location"`
	ProcessedAt  int64  `json:"processed_at"`
}

func main() {
	brokers := []string{"localhost:9092"} // Kafka broker
	consumerGroup := "clickstream-processor"
	rawTopic := "raw-events"
	processedTopic := "processed-events"

	// Configure Sarama
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create consumer group
	client, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer client.Close()

	// Create producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle system signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigchan
		cancel()
	}()

	handler := &ProcessorHandler{producer: producer, processedTopic: processedTopic}

	log.Println("Processor service started. Listening for raw events...")

	for {
		if err := client.Consume(ctx, []string{rawTopic}, handler); err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

// ProcessorHandler implements sarama.ConsumerGroupHandler
type ProcessorHandler struct {
	producer       sarama.SyncProducer
	processedTopic string
}

func (h *ProcessorHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ProcessorHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ProcessorHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var raw RawEvent
		if err := json.Unmarshal(msg.Value, &raw); err != nil {
			log.Printf("Failed to unmarshal raw event: %v", err)
			continue
		}

		// Example enrichment: add dummy device + geolocation info
		processed := ProcessedEvent{
			UserID:      raw.UserID,
			EventType:   raw.EventType,
			Timestamp:   raw.Timestamp,
			Device:      "mobile",
			GeoLocation: "India",
			ProcessedAt: msg.Timestamp.Unix(),
		}

		data, err := json.Marshal(processed)
		if err != nil {
			log.Printf("Failed to marshal processed event: %v", err)
			continue
		}

		// Publish to processed-events topic
		_, _, err = h.producer.SendMessage(&sarama.ProducerMessage{
			Topic: h.processedTopic,
			Value: sarama.ByteEncoder(data),
		})
		if err != nil {
			log.Printf("Failed to send processed event: %v", err)
			continue
		}

		log.Printf("Processed event for user %s", raw.UserID)
		session.MarkMessage(msg, "")
	}
	return nil
}
