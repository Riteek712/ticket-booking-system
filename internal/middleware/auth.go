package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// JWTProtected defines the JWT middleware configuration
func JWTProtected() fiber.Handler {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		panic("JWT_SECRET is not set in environment variables")
	}

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(secretKey),
		ContextKey:   "jwt",           // Store token payload in this context key
		ErrorHandler: jwtErrorHandler, // Custom error handler
	})
}

// jwtErrorHandler handles JWT authentication errors
func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   "Unauthorized: " + err.Error(),
	})
}
