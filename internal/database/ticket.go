package database

import "gorm.io/gorm"

// Ticket represents a ticket for an event.
type Ticket struct {
	gorm.Model `swaggerignore:"true"`
	Email      string `gorm:"type:varchar(255);not null"`        // Email of the ticket holder
	TicketID   string `gorm:"type:varchar(255);unique;not null"` // Unique ticket ID
	EventID    string `gorm:"not null" json:"event_id"`          // ID of the event associated with the ticket
	Quantity   int    `gorm:"not null" json:"quantity"`          // Number of tickets booked
	Event      Event  `gorm:"constraint:OnDelete:CASCADE;"`      // Relationship with Event
}

// TicketBookingReq represents the request payload for booking a ticket.
type TicketBookingReq struct {
	Email    string `json:"email" validate:"required,email"`    // Email of the ticket holder
	EventID  string `json:"event_id" validate:"required"`       // ID of the event to book
	Quantity int    `json:"quantity" validate:"required,min=1"` // Number of tickets to book
}
