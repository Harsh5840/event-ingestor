package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/go-event-streaming/internal/processor"
)

func main () {
	fmt.Println("Starting processor Service")

	ctx, cancel := context.WithCancel(context.Background())    //here , we declared the context 
	defer cancel()

	//graceful shutdown
	sigs := make(chan os.Signal , 1) 
	signal.Notify(sigs, syscall.SIGINT , syscall.SIGTERM)  //here 

	go func ()  {
		<-sigs
		fmt.Println("Shutting down processor service")
		cancel()
	}()

	//starting processor logic
	if  err := processor.Start(ctx); err != nil {
		log.Fatalf("Processor service error : %v" , err)
	}

	fmt.Println("Processor Service Stopped")
}