package server

import (
	"log"
	"strings"
	"ticketing/internal/database"
	"ticketing/internal/utils"

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
// @Security BearerAuth
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

	// Extract the user_id from the token
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
	}

	// Remove the 'Bearer ' prefix from the token string
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Create a new event instance with a generated UUID
	event := &database.Event{
		EventID:     uuid.New().String(), // Assuming you have a UUID field in your Event struct
		Name:        dto.Name,
		Description: dto.Description,
		Capacity:    dto.Capacity,
		UserID:      userID,
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
// @Security BearerAuth
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

// updateEvent updates an existing event.
// @Summary Update an event
// @Description Update an event's name, description, and capacity by event ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body database.CreateEventDTO true "Updated Event Data"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [put]
// @Security BearerAuth
func (s *FiberServer) updateEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	var dto database.CreateEventDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Fetch the event from the database
	event, err := s.db.GetEvent(eventID)
	if err != nil {
		if err == database.ErrEventNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve event"})
	}

	// Update the event fields
	event.Name = dto.Name
	event.Description = dto.Description
	event.Capacity = dto.Capacity

	// Save the updated event to the database
	if err := s.db.UpdateEvent(event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update event"})
	}

	return c.JSON(event)
}

// deleteEvent deletes an event by ID.
// @Summary Delete an event
// @Description Delete an event by event ID
// @Tags events
// @Param id path string true "Event ID"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [delete]
// @Security BearerAuth
func (s *FiberServer) deleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	// Try to delete the event from the database
	if err := s.db.DeleteEvent(eventID); err != nil {
		if err == database.ErrEventNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete event"})
	}

	// Return a 204 No Content response if the deletion was successful
	return c.SendStatus(fiber.StatusNoContent)
}
