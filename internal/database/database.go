package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	CreateEvent(event *Event) error
	GetEvent(uniqueID string) (*Event, error)
	UpdateEvent(event *Event) error
	DeleteEvent(uniqueID, userId string) error
	GetTotalTicketsSold(eventID string) (int, error)
	CreateTicket(ticket *Ticket) error
	GetTicket(ticketID string) (*Ticket, error)
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}

type service struct {
	db *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

// New initializes the database connection using GORM.
func New() (Service, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		host, username, password, database, port, schema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Disable foreign key checks
	db.Exec("SET CONSTRAINTS ALL DEFERRED;")

	// Ensure the correct order of migration
	if err := db.AutoMigrate(&Event{}, &Ticket{}, &User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance, nil
}

// Health checks the health of the database connection.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Check if the database is reachable
	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	// Ping the database
	if err := sqlDB.PingContext(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	// Database is up
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats
}

// Close closes the database connection.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Defined the error for event not found
var ErrEventNotFound = errors.New("event not found")

// CreateEvent creates a new event in the database.
func (s *service) CreateEvent(event *Event) error {
	return s.db.Create(event).Error
}

// GetEvent retrieves an event by its unique ID.
func (s *service) GetEvent(eventID string) (*Event, error) {
	var event Event
	if err := s.db.First(&event, "event_id = ?", eventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return &event, nil
}

// UpdateEvent updates an existing event in the database.
func (s *service) UpdateEvent(event *Event) error {
	return s.db.Save(event).Error
}

// DeleteEvent deletes an event by its unique ID.
func (s *service) DeleteEvent(uniqueID, userID string) error {
	return s.db.Delete(&Event{}, "unique_id = ? AND user_id = ?", uniqueID, userID).Error
}

// GetTotalTicketsSold retrieves the total number of tickets sold for a specific event.
func (s *service) GetTotalTicketsSold(eventID string) (int, error) {
	var totalSold int64
	if err := s.db.Model(&Ticket{}).
		Where("event_id = ?", eventID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalSold).Error; err != nil {
		return 0, err
	}
	return int(totalSold), nil
}

// CreateTicket saves a new ticket in the database.
func (s *service) CreateTicket(ticket *Ticket) error {
	return s.db.Model(&Ticket{}).Create(ticket).Error
}
func (s *service) GetTicket(ticketID string) (*Ticket, error) {
	var ticket Ticket
	if err := s.db.First(&ticket, "ticket_id = ?", ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}
	return &ticket, nil
}

func (s *service) CreateUser(user *User) error {
	return s.db.Model(&User{}).Create(user).Error
}

func (s *service) GetUserByEmail(email string) (*User, error) {
	var res User
	if err := s.db.Model(&User{}).Where("email = ?", email).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
