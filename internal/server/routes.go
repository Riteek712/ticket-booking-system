package server

import (
	"ticketing/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/swagger/*", swagger.HandlerDefault) // Default serves swagger at /swagger/index.html

	s.App.Post("/users/register", s.RegisterUser)
	s.App.Post("/users/login", s.LoginUser)

	s.App.Post("/events", middleware.JWTProtected(), s.createEvent)
	s.App.Get("/events/:id", middleware.JWTProtected(), s.getEvent)
	s.App.Put("/events/:id", middleware.JWTProtected(), s.updateEvent)
	s.App.Delete("/events/:id", middleware.JWTProtected(), s.deleteEvent)

	s.App.Post("/tickets", middleware.JWTProtected(), s.bookTicket)

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
