package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/swagger/*", swagger.HandlerDefault) // Default serves swagger at /swagger/index.html

	s.App.Post("/events", s.createEvent)
	s.App.Get("/events/:id", s.getEvent)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
