package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-event-streaming/internal/api"
)

func main() {
	fmt.Println("Starting API service")

	ctx, stop := signal.NotifyContext(context.Background() , syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Run server in a go routine
	go func ()  {
		if err := server.ListenAndServe(); err!= nil && err!= http.ErrServerClosed {
			log.Fatalf("API server failed: %v", err)
		}
	}()
}