package server

import (
	"log"
	"ticketing/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// bookTicket handles ticket booking requests.
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
// @Router /tickets [post]
func (s *FiberServer) bookTicket(c *fiber.Ctx) error {
	var req database.TicketBookingReq

	// Parse the JSON body into the request struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Check if the event exists
	event, err := s.db.GetEvent(req.EventID) // Assuming you have a GetEventByID function
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	// Calculate total sold tickets
	totalSold, err := s.db.GetTotalTicketsSold(req.EventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not calculate total sold tickets"})
	}

	// Calculate remaining capacity
	remainingCapacity := event.Capacity - totalSold

	// Check if the requested capacity exceeds the remaining capacity
	if req.Quantity > remainingCapacity {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient capacity for this event"})
	}

	// Create a new ticket instance
	ticket := &database.Ticket{
		Email:    req.Email,
		TicketID: uuid.New().String(), // Generate a unique TicketID
		EventID:  req.EventID,
		Quantity: req.Quantity,
	}

	// Call the CreateTicket method from the database service
	if err := s.db.CreateTicket(ticket); err != nil {
		log.Printf("Error creating ticket: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create ticket"})
	}

	// Respond with the created ticket
	return c.Status(fiber.StatusCreated).JSON(ticket)
}
