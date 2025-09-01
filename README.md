# Event Streaming System (Go + Kafka)

## 📌 Overview
This project is a distributed event streaming platform built with **Go** and **Kafka**.  
It demonstrates the ability to handle **large incoming data streams** by simulating ingestion, processing, and serving APIs.

## 🏗 Architecture
- **Ingestor Service** → Reads raw events and pushes to Kafka.
- **Processor Service** → Consumes events, transforms them, and republishes them.
- **API Service** → Exposes REST endpoints to query processed data.

## 🔧 Tech Stack
- Go 1.22
- Kafka (via `segmentio/kafka-go`)
- Docker & Docker Compose
- Kubernetes (for optional deployment)
- Logrus for logging

## ⚡ Getting Started
### 1. Clone repo
```bash
git clone https://github.com/your-username/event-streaming-system.git
cd event-streaming-system
