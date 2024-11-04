package server

import "github.com/gofiber/fiber/v2"

func (s *FiberServer) ticket(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
