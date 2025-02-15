package broker

import (
	"encoding/json"
	"log"
	"ticketing/internal/database"

	"github.com/google/uuid"
)

func ConsumeBookingRequests(db database.Service) {
	ch := GetRabbitMQChannel()
	msgs, err := ch.Consume(
		"ticket_booking_queue",
		"",    // Consumer tag
		true,  // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var req database.TicketBookingReq
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Process booking
		err = processBooking(db, req)
		if err != nil {
			log.Printf("Failed to book ticket: %v", err)
		} else {
			log.Println("Ticket booked successfully!")
		}
	}
}

func processBooking(db database.Service, req database.TicketBookingReq) error {
	// Check if event exists
	event, err := db.GetEvent(req.EventID)
	if err != nil {
		log.Printf("Event not found: %v", err)
		return err
	}

	// Check ticket availability
	totalSold, err := db.GetTotalTicketsSold(req.EventID)
	if err != nil {
		log.Printf("Error fetching sold tickets: %v", err)
		return err
	}

	remainingCapacity := event.Capacity - totalSold
	if req.Quantity > remainingCapacity {
		log.Printf("Insufficient capacity for event %v", req.EventID)
		return err
	}

	// Create ticket
	ticket := &database.Ticket{
		TicketID: uuid.New().String(),
		EventID:  req.EventID,
		Email:    req.Email,
		Quantity: req.Quantity,
	}

	return db.CreateTicket(ticket)
}
