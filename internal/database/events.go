package database

import (
	"gorm.io/gorm"
)

// Event represents the event model in the database.// Event represents the event model in the database.
type Event struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string   `gorm:"type:varchar(255);not null"`        // Name of the event
	Description string   `gorm:"type:text;not null"`                // Event description
	EventID     string   `gorm:"type:varchar(255);unique;not null"` // Unique event identifier
	Capacity    int      `gorm:"not null" json:"capacity"`          // Total capacity of the event
	UserID      string   ` json:"user_id"`
	Tickets     []Ticket `gorm:"foreignKey:EventID" json:"tickets"` // Associated tickets
}

// CreateEventDTO represents the data for creating a new event.
type CreateEventDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=255"` // Event name
	Description string `json:"description" validate:"required"`        // Event description
	Capacity    int    `json:"capacity" validate:"required,min=1"`     // Total capacity of the event
}
