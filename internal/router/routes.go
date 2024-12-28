package router

import (
	"ticketing/internal/database"
	"ticketing/internal/handler"
	"ticketing/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db database.Service) {
	// Handlers
	helloHandler := handler.NewHelloHandler(db)
	userHandler := handler.NewUserHandler(db)
	eventHandler := handler.NewEventHandler(db)
	ticketHandler := handler.NewTicketHandler(db)

	rateLimit := middleware.RateLimitMiddleware(5, 10*time.Second)

	// Routes
	app.Get("/", helloHandler.HelloWorld)
	app.Get("/health", helloHandler.Health)

	app.Post("/users/register", userHandler.RegisterUser)
	app.Post("/users/login", userHandler.LoginUser)

	app.Post("/events", middleware.JWTProtected(), rateLimit, eventHandler.CreateEvent)
	app.Get("/events/:id", middleware.JWTProtected(), rateLimit, eventHandler.GetEvent)
	app.Put("/events/:id", middleware.JWTProtected(), rateLimit, eventHandler.UpdateEvent)
	app.Delete("/events/:id", middleware.JWTProtected(), rateLimit, eventHandler.DeleteEvent)

	app.Post("/tickets", middleware.JWTProtected(), rateLimit, ticketHandler.BookTicket)
}
