package database

import "time"

// User represents a user in the ticket-event management system.
type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        string    `gorm:"not null" json:"user_id"`
	FirstName     string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName      string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Email         string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash  string    `gorm:"type:text;not null" json:"-"` // Avoid sending the password hash in JSON responses
	Phone         string    `gorm:"type:varchar(15)" json:"phone"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TicketsBooked []Ticket  `gorm:"foreignKey:UserID" json:"tickets_booked"`
}

type SignUpDTO struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Phone     string `json:"phoneNumber"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
