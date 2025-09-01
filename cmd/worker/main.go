package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"clickstreamx/internal/worker"
)

func main() {
	fmt.Println("🚀 Starting Worker Service...")

	// Create a new worker that consumes from the message queue
	w := worker.NewWorker()

	// Run the worker in a goroutine
	go func() {
		if err := w.Start(); err != nil {
			log.Fatalf("❌ Worker failed: %v", err)
		}
	}()

	// Graceful shutdown on interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("🛑 Worker shutting down...")
	w.Stop()
}
