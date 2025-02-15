package broker

import (
	"encoding/json"
	"log"
	"ticketing/internal/database"

	"github.com/streadway/amqp"
)

func PublishBookingRequest(ticketReq database.TicketBookingReq) error {
	ch := GetRabbitMQChannel()

	body, err := json.Marshal(ticketReq)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",                     // Exchange
		"ticket_booking_queue", // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Println("Booking request published to queue")
	return nil
}
