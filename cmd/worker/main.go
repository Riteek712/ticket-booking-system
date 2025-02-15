package main

import (
	"log"
	"ticketing/internal/broker"
	"ticketing/internal/database"
)

func main() {
	// Initialize database
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize RabbitMQ
	_, _, err = broker.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	log.Println("Worker started, listening for ticket booking requests...")
	broker.ConsumeBookingRequests(db)
}
