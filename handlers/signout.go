package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) Signout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "x-access-token"
	cookie.Value = ""
	cookie.Domain = "localhost"
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = "None"
	cookie.Expires = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c.Cookie(cookie)

	return c.SendStatus(204)
}
