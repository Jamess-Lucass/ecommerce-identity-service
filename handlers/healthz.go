package handlers

import "github.com/gofiber/fiber/v2"

func (s *Server) Healthz(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).Send([]byte("Healthy"))
}
