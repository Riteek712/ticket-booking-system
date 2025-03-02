package main

import (
	"log"
	"ticketing/internal/broker"
	"ticketing/internal/database"
	"ticketing/internal/queue"
	"ticketing/internal/router"
	"ticketing/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"sync"
)

// Starts API server
func startAPI(wg *sync.WaitGroup) {
	defer wg.Done()
	app := fiber.New()
	utils.InitRedis()

	db, err := database.New()
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	queueService, err := queue.NewQueueService()
	if err != nil {
		log.Fatalf("could not init the queue service: %v", err)
	}

	app.Use(cors.New())
	router.RegisterRoutes(app, db, *queueService)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(app.Listen(":3000"))
}

// Starts Worker
func startWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	_, _, err = broker.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	log.Println("Worker started, listening for ticket booking requests...")
	broker.ConsumeBookingRequests(db)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go startAPI(&wg)
	go startWorker(&wg)

	wg.Wait()
}
