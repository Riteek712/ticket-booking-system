package main

import (
	"fmt"
	"os"
	"strconv"
	"ticketing/internal/server"

	_ "ticketing/docs" // Import Swagger docs here

	_ "github.com/joho/godotenv/autoload"
)

// @title Ticket-Booking API
// @version 1.0
// @description This is a sample Ticket-Booking API  server for a Fiber app.
// @host 127.0.0.1:8080
// @BasePath /
func main() {

	server := server.New()

	// @securityDefinitions.apiKey BearerAuth
	// @in header
	// @name Authorization
	server.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
