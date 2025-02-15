package queue

import (
	"encoding/json"
	"log"
	"ticketing/internal/database"

	"github.com/streadway/amqp"
)

// Service handles RabbitMQ operations
type Service struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewQueueService initializes a RabbitMQ connection and declares a queue
func NewQueueService() (*Service, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"ticket_booking_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Service{conn: conn, channel: ch, queue: q}, nil
}

// PublishTicketRequest enqueues a ticket booking request
func (s *Service) PublishTicketRequest(req database.TicketBookingReq) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = s.channel.Publish(
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	return nil
}

// GetQueueLength returns the length of the queue
func (s *Service) GetQueueLength(eventID string) (int, error) {
	q, err := s.channel.QueueInspect(s.queue.Name)
	if err != nil {
		return 0, err
	}
	return q.Messages, nil
}

// Close closes the RabbitMQ connection
func (s *Service) Close() {
	s.channel.Close()
	s.conn.Close()
}
