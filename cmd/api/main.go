package main

import (
	"log"
	"ticketing/internal/database"
	"ticketing/internal/queue"
	"ticketing/internal/router"
	"ticketing/internal/utils"

	_ "ticketing/docs" // Import the generated docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger" // Swagger UI middleware
)

// @title Ticket-Booking API
// @version 1.0
// @description This is a sample Ticket-Booking API server for a Fiber app.
// @host 127.0.0.1:3000
// @BasePath /
func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Initialize Redis
	utils.InitRedis()

	// Set up the database connection
	db, err := database.New() // Assuming you have a function that sets up the database
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	// Set up the queue service
	queueService, err := queue.NewQueueService()
	if err != nil {
		log.Fatalf("could not init the queue service: %v", err)
	}
	// Middleware
	app.Use(cors.New()) // Enable CORS if required

	// Register routes
	// @securityDefinitions.apiKey BearerAuth
	// @in header
	// @name Authorization
	router.RegisterRoutes(app, db, *queueService)

	// Change Swagger UI route to avoid conflict (for example, /api-docs)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
