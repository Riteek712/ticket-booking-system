package router

import (
	"ticketing/internal/database"
	"ticketing/internal/handler"
	"ticketing/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db database.Service) {
	// Handlers
	helloHandler := handler.NewHelloHandler(db)
	userHandler := handler.NewUserHandler(db)
	eventHandler := handler.NewEventHandler(db)
	ticketHandler := handler.NewTicketHandler(db)

	// Routes
	app.Get("/", helloHandler.HelloWorld)
	app.Get("/health", helloHandler.Health)

	app.Post("/users/register", userHandler.RegisterUser)
	app.Post("/users/login", userHandler.LoginUser)

	app.Post("/events", middleware.JWTProtected(), eventHandler.CreateEvent)
	app.Get("/events/:id", middleware.JWTProtected(), eventHandler.GetEvent)
	app.Put("/events/:id", middleware.JWTProtected(), eventHandler.UpdateEvent)
	app.Delete("/events/:id", middleware.JWTProtected(), eventHandler.DeleteEvent)

	app.Post("/tickets", middleware.JWTProtected(), ticketHandler.BookTicket)
}
