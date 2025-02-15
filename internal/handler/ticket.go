package handler

import (
	"log"
	"ticketing/internal/database"
	"ticketing/internal/queue"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TicketHandler represents the ticket-related HTTP handlers
// @BasePath /api/v1

type TicketHandler struct {
	db    database.Service
	queue queue.Service // RabbitMQ service
}

// NewTicketHandler creates a new instance of TicketHandler
func NewTicketHandler(db database.Service, queue queue.Service) *TicketHandler {
	return &TicketHandler{db: db, queue: queue}
}

// GetQueueLength returns the number of pending ticket requests for an event
// @Summary Get event queue length
// @Description Returns the number of pending ticket requests for a specific event
// @Tags Tickets
// @Accept  json
// @Produce  json
// @Param eventID path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /queue/{eventID}/length [get]
func (h *TicketHandler) GetQueueLength(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	queueLength, err := h.queue.GetQueueLength(eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch queue length"})
	}

	return c.JSON(fiber.Map{"event_id": eventID, "queue_length": queueLength})
}

// GetTicketDetails fetches a ticket and the associated event information
// @Summary Get ticket details
// @Description Fetches details of a specific ticket along with the associated event information
// @Tags Tickets
// @Accept  json
// @Produce  json
// @Param ticketID path string true "Ticket ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ticket/{ticketID} [get]
func (h *TicketHandler) GetTicketDetails(c *fiber.Ctx) error {
	ticketID := c.Params("ticketID")

	// Fetch ticket details
	ticket, err := h.db.GetTicket(ticketID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	}

	// Fetch associated event details
	event, err := h.db.GetEvent(ticket.EventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Event details not found"})
	}

	return c.JSON(fiber.Map{"ticket": ticket, "event": event})
}

// AddTicketToQueue enqueues a ticket booking request to RabbitMQ
// @Summary Add ticket booking request to queue
// @Description Enqueues a ticket booking request for an event in RabbitMQ
// @Tags Tickets
// @Accept  json
// @Produce  json
// @Param request body database.TicketBookingReq true "Ticket booking request payload"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ticket/book [post]
func (h *TicketHandler) AddTicketToQueue(c *fiber.Ctx) error {
	var req database.TicketBookingReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Generate a unique TicketID
	req.TicketID = uuid.New().String()

	// Publish the request to RabbitMQ
	if err := h.queue.PublishTicketRequest(req); err != nil {
		log.Printf("Failed to enqueue ticket request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not enqueue ticket request"})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Ticket booking request added to queue", "ticket_id": req.TicketID})
}
