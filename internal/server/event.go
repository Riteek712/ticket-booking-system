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

// getEvent retrieves an event by ID.
// @Summary Get an event by ID
// @Description Retrieve an event's details by its unique event ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} database.Event
// @Failure 404 {object} map[string]interface{} "Event not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /events/{id} [get]
func (s *FiberServer) getEvent(c *fiber.Ctx) error {
	eventID := c.Params("id") // Get the event ID from the URL path

	// Retrieve the event from the database
	event, err := s.db.GetEvent(eventID)
	if err != nil {
		log.Printf("Error retrieving event: %v", err)
		if err == database.ErrEventNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve event"})
	}

	// Respond with the found event
	return c.JSON(event)
}
