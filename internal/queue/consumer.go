package queue

import (
	"encoding/json"
	"log"
	"ticketing/internal/database"
)

// ProcessQueue consumes messages from RabbitMQ and books tickets
func (s *Service) ProcessQueue(db database.Service) {
	msgs, err := s.channel.Consume(
		s.queue.Name,
		"",
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	go func() {
		for d := range msgs {
			var req database.TicketBookingReq
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Process the booking request
			ticket := &database.Ticket{
				TicketID: req.TicketID,
				EventID:  req.EventID,
				Email:    req.Email,
				Quantity: req.Quantity,
			}

			if err := db.CreateTicket(ticket); err != nil {
				log.Printf("Failed to book ticket: %v", err)
				continue
			}

			log.Printf("Ticket booked successfully: %+v", ticket)
		}
	}()
}
