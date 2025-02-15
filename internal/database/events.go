package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// Event represents the event model in the database.// Event represents the event model in the database.
type Event struct {
	gorm.Model   `swaggerignore:"true"`
	Name         string             `gorm:"type:varchar(255);not null"`        // Name of the event
	Description  string             `gorm:"type:text;not null"`                // Event description
	EventID      string             `gorm:"type:varchar(255);unique;not null"` // Unique event identifier
	Capacity     int                `gorm:"not null" json:"capacity"`          // Total capacity of the event
	UserID       string             `json:"user_id"`
	EventDetails EventDetailsStruct ` gorm:"type:jsonb" json:"event_details"` // Additional event details (not stored in DB)
	// Tickets      []Ticket           `gorm:"foreignKey:EventID" json:"tickets"` // Associated tickets
}

// CreateEventDTO represents the data for creating a new event.
type CreateEventDTO struct {
	Name         string             `json:"name" validate:"required,min=3,max=255"` // Event name
	Description  string             `json:"description" validate:"required"`        // Event description
	Capacity     int                `json:"capacity" validate:"required,min=1"`     // Total capacity of the event
	EventDetails EventDetailsStruct ` gorm:"type:jsonb" json:"event_details"`
}

type EventDetailsStruct struct {
	Details map[string]interface{} ` gorm:"type:jsonb" json:"details"`
}

// Scan implements the `sql.Scanner` interface for EventDetails.
func (m EventDetailsStruct) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if data, ok := value.([]byte); ok {
		err := json.Unmarshal(data, &m) // Unmarshal into the map
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid data type for EventDetails")
}

// Value implements the `driver.Valuer` interface for EventDetails.
func (m EventDetailsStruct) Value() (driver.Value, error) {
	serialized, err := json.Marshal(m) // Marshal map into JSON
	if err != nil {
		return nil, err
	}
	return serialized, nil // Return the serialized JSON
}
