package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/swagger/*", swagger.HandlerDefault) // Default serves swagger at /swagger/index.html

	s.App.Post("/users/register", s.RegisterUser)
	s.App.Post("users/login", s.LoginUser)

	s.App.Post("/events", s.createEvent)
	s.App.Get("/events/:id", s.getEvent)
	s.App.Put("/events/:id", s.updateEvent)
	s.App.Delete("/events/:id", s.deleteEvent)

	s.App.Post("/tickets", s.bookTicket)

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
