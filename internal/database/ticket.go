package database

import "gorm.io/gorm"

// Ticket represents a ticket for an event.
type Ticket struct {
	gorm.Model `swaggerignore:"true"` // Ignore the model fields in Swagger documentation
	Email      string                 `gorm:"type:varchar(255);not null"`        // Email of the ticket holder
	TicketID   string                 `gorm:"type:varchar(255);unique;not null"` // Unique ticket ID
	EventID    string                 `gorm:"type:varchar(255);not null"`        // ID of the event associated with the ticket
	Capacity   int                    `gorm:"not null"`                          // Number of tickets booked
}

// TicketBookingReq represents the request payload for booking a ticket.
type TicketBookingReq struct {
	Email    string `json:"email" validate:"required"`    // Email of the ticket holder
	EventID  string `json:"event_id" validate:"required"` // ID of the event to book
	Capacity int    `json:"capacity" validate:"required"` // Number of tickets to book
}
