package handler

import (
	"ticketing/internal/database"

	"github.com/gofiber/fiber/v2"
)

// HealthResponse defines the structure of the health check response.
type HealthResponse struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details"`
}

// HelloHandler represents the handler for the Hello world and health check endpoints.
type HelloHandler struct {
	db database.Service
}

// NewHelloHandler creates a new instance of HelloHandler with the database service.
func NewHelloHandler(db database.Service) *HelloHandler {
	return &HelloHandler{db: db}
}

func (h *HelloHandler) HelloWorld(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}
	return c.JSON(resp)
}

func (h *HelloHandler) Health(c *fiber.Ctx) error {
	// Prepare the response in the expected format
	health := h.db.Health()
	// Return the response in HealthResponse format
	resp := HealthResponse{
		Status: health["status"],
		Details: map[string]string{
			"database": health["message"], // You might want to adjust these fields
		},
	}
	return c.JSON(resp)
}
