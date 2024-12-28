package middleware

import (
	"context"
	"fmt"
	"ticketing/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimitMiddleware applies rate limiting using Redis.
func RateLimitMiddleware(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Unique key for the user (e.g., based on IP or user ID)
		userID := c.IP() // Or use c.Locals("userID") for authenticated users
		redisKey := fmt.Sprintf("rate_limit:%s", userID)

		// Increment request count
		count, err := utils.Rdb.Incr(context.Background(), redisKey).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		// Set TTL on the key if it's the first request
		if count == 1 {
			_, err := utils.Rdb.Expire(context.Background(), redisKey, window).Result()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}
		}

		// Check if the request count exceeds the limit
		if count > int64(limit) {
			remainingTTL, _ := utils.Rdb.TTL(context.Background(), redisKey).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message":    "Rate limit exceeded",
				"retryAfter": remainingTTL.Seconds(),
			})
		}

		// Proceed to the next handler
		return c.Next()
	}
}
