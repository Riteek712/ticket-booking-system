package router

import (
	"ticketing/internal/database"
	"ticketing/internal/handler"
	"ticketing/internal/middleware"
	"ticketing/internal/queue"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db database.Service, queueService queue.Service) {
	// Handlers
	helloHandler := handler.NewHelloHandler(db)
	userHandler := handler.NewUserHandler(db)
	eventHandler := handler.NewEventHandler(db)
	ticketHandler := handler.NewTicketHandler(db, queueService)

	rateLimit := middleware.RateLimitMiddleware(5, 5*time.Second)

	// Routes
	app.Get("/", helloHandler.HelloWorld)
	app.Get("/health", helloHandler.Health)

	app.Post("/users/register", userHandler.RegisterUser)
	app.Post("/users/login", userHandler.LoginUser)

	app.Post("/events", middleware.JWTProtected(), rateLimit, eventHandler.CreateEvent)
	app.Get("/events/:id", middleware.JWTProtected(), rateLimit, eventHandler.GetEvent)
	app.Put("/events/:id", middleware.JWTProtected(), rateLimit, eventHandler.UpdateEvent)
	app.Delete("/events/:id", middleware.JWTProtected(), rateLimit, eventHandler.DeleteEvent)

	app.Post("/tickets", rateLimit, ticketHandler.AddTicketToQueue)
	app.Get("/tickets/:ticketID", rateLimit, ticketHandler.GetTicketDetails)
	app.Get("/queue/:eventID/length", rateLimit, ticketHandler.GetQueueLength)
}
