package database

import "time"

// User represents a user in the ticket-event management system.
type User struct {
	UserID       string    `gorm:"primaryKey;not null" json:"user_id"`
	FirstName    string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName     string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	Phone        string    `gorm:"type:varchar(15)" json:"phone"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	// TicketsBooked []Ticket  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"tickets_booked"`
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
