package handlers

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) Signout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "x-access-token"
	cookie.Value = ""
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = "None"
	cookie.Expires = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	if strings.EqualFold(os.Getenv("ENVIRONMENT"), "production") {
		cookie.Domain = "jameslucas.uk"
	} else {
		cookie.Domain = "localhost"
	}

	c.Cookie(cookie)

	return c.SendStatus(fiber.StatusNoContent)
}
