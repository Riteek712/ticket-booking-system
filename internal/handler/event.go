package handler

import (
	"log"
	"strings"
	"ticketing/internal/database"
	"ticketing/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	DB database.Service
}

func NewEventHandler(db database.Service) *EventHandler {
	return &EventHandler{DB: db}
}

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
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var dto database.CreateEventDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if dto.Name == "" || dto.Description == "" || dto.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Description, and Capacity are required"})
	}

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	event := &database.Event{
		EventID:      uuid.New().String(),
		Name:         dto.Name,
		Description:  dto.Description,
		Capacity:     dto.Capacity,
		UserID:       userID,
		EventDetails: dto.EventDetails,
	}

	if err := h.DB.CreateEvent(event); err != nil {
		log.Printf("Error creating event: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create event"})
	}

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
func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	event, err := h.DB.GetEvent(eventID)
	if err != nil {
		if err == database.ErrEventNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve event"})
	}

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
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	var dto database.CreateEventDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	event, err := h.DB.GetEvent(eventID)
	if err != nil {
		if err == database.ErrEventNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve event"})
	}

	if event.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Your are not allowed to edit the datafor this event."})
	}

	event.Name = dto.Name
	event.Description = dto.Description
	event.Capacity = dto.Capacity

	if err := h.DB.UpdateEvent(event); err != nil {
		log.Printf("Error updating event: %v", err)
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
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	userID, err := utils.ExtractUserID(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	if err := h.DB.DeleteEvent(eventID, userID); err != nil {
		if err == database.ErrEventNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete event"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
