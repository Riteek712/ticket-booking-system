package handler

import (
	"log"
	"ticketing/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TicketHandler represents the ticket-related HTTP handlers
type TicketHandler struct {
	db database.Service // Assume DatabaseService interface manages database operations
}

// NewTicketHandler creates a new instance of TicketHandler
func NewTicketHandler(db database.Service) *TicketHandler {
	return &TicketHandler{db: db}
}

// BookTicket handles ticket booking requests.
// @Summary Book a ticket
// @Description Book a ticket for an event
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body database.TicketBookingReq true "Ticket Booking Request"
// @Success 201 {object} database.Ticket
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /tickets [post]
func (h *TicketHandler) BookTicket(c *fiber.Ctx) error {
	var req database.TicketBookingReq

	// Parse the JSON body into the request struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Check if the event exists
	event, err := h.db.GetEvent(req.EventID)
	if err != nil {
		log.Printf("Error fetching event: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	// Calculate total sold tickets
	totalSold, err := h.db.GetTotalTicketsSold(req.EventID)
	if err != nil {
		log.Printf("Error fetching sold tickets: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not calculate total sold tickets"})
	}

	// Calculate remaining capacity
	remainingCapacity := event.Capacity - totalSold

	// Check if the requested quantity exceeds the remaining capacity
	if req.Quantity > remainingCapacity {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient capacity for this event"})
	}

	// Create a new ticket instance
	ticket := &database.Ticket{
		TicketID: uuid.New().String(), // Generate a unique TicketID
		EventID:  req.EventID,
		Email:    req.Email,
		Quantity: req.Quantity,
	}

	// Save the ticket to the database
	if err := h.db.CreateTicket(ticket); err != nil {
		log.Printf("Error creating ticket: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create ticket"})
	}

	// Respond with the created ticket
	return c.Status(fiber.StatusCreated).JSON(ticket)
}
