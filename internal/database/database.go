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
	DeleteEvent(uniqueID string) error
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
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		host, username, password, database, port, schema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Event{}); err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
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
		return nil, err
	}
	return &event, nil
}

// UpdateEvent updates an existing event in the database.
func (s *service) UpdateEvent(event *Event) error {
	return s.db.Save(event).Error
}

// DeleteEvent deletes an event by its unique ID.
func (s *service) DeleteEvent(uniqueID string) error {
	return s.db.Delete(&Event{}, "unique_id = ?", uniqueID).Error
}
