package database

import "gorm.io/gorm"

// Ticket represents a ticket for an event.
type Ticket struct {
	gorm.Model `swaggerignore:"true"`
	Email      string `gorm:"type:varchar(255);not null"`
	TicketID   string `gorm:"type:varchar(255);unique;not null"`
	EventID    string `gorm:"not null" json:"event_id"`
	UserID     string `gorm:"not null" json:"user_id"` // Ensure this is string
	Quantity   int    `gorm:"not null" json:"quantity"`
	// Event      Event  `gorm:"constraint:OnDelete:CASCADE;"`
	// User       User   `gorm:"foreignKey:UserID;references:UserID"` // Explicitly reference UserID
}

// TicketBookingReq represents the request payload for booking a ticket.
type TicketBookingReq struct {
	TicketID string `json:"ticket_id"`
	Email    string `json:"email" validate:"required,email"`    // Email of the ticket holder
	EventID  string `json:"event_id" validate:"required"`       // ID of the event to book
	Quantity int    `json:"quantity" validate:"required,min=1"` // Number of tickets to book
}
