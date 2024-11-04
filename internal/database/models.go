package database

import (
	"gorm.io/gorm"
)

// Event represents the event model in the database.
type Event struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;not null"`
	EventID     string `gorm:"type:varchar(255);unique;not null"`
	Capacity    int    `gorm:"not null"`
}
