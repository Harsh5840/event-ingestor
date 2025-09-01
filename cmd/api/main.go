package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/segmentio/kafka-go"
)

type ClickEvent struct {
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	PageURL   string `json:"page_url"`
	Action    string `json:"action"`
}

var kafkaWriter *kafka.Writer

func main() {
	// Kafka setup
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "clickstream_events")

	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	http.HandleFunc("/collect", collectHandler)

	log.Println("Clickstream API running on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("server error:", err)
	}
}

func collectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var event ClickEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Marshal to JSON
	msg, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "serialization error", http.StatusInternalServerError)
		return
	}

	// Write to Kafka
	err = kafkaWriter.WriteMessages(r.Context(),
		kafka.Message{Value: msg},
	)
	if err != nil {
		http.Error(w, "failed to push event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("event received"))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
