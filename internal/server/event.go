package server

import (
	"log"
	"ticketing/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// createEvent creates a new event.
// @Summary Create a new event
// @Description Create a new event with name, description, and capacity
// @Tags events
// @Accept json
// @Produce json
// @Param event body database.CreateEventDTO true "Event Data"
// @Success 201 {object} database.Event
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events [post]
func (s *FiberServer) createEvent(c *fiber.Ctx) error {
	var dto database.CreateEventDTO

	// Parse the JSON body into the DTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	// Validate the DTO (optional, you can use a validation library)
	if dto.Name == "" || dto.Description == "" || dto.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Description, and Capacity are required"})
	}

	// Create a new event instance with a generated UUID
	event := &database.Event{
		EventID:     uuid.New().String(), // Assuming you have a UUID field in your Event struct
		Name:        dto.Name,
		Description: dto.Description,
		Capacity:    dto.Capacity,
	}

	// Call the CreateEvent method from the database service
	if err := s.db.CreateEvent(event); err != nil {
		log.Printf("Error creating event: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create event"})
	}

	// Respond with the created event
	return c.Status(fiber.StatusCreated).JSON(event)
}
