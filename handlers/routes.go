package handlers

import (
	"github.com/Jamess-Lucass/ecommerce-identity-service/middleware"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *Server) Start() error {
	f := fiber.New()
	f.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		MaxAge:           0,
	}))

	f.Use(fiberzap.New(fiberzap.Config{
		Logger: s.logger,
	}))

	f.Get("/api/v1/oauth/authorize/google", s.RedirectGoogleAuthorize)
	f.Get("/api/v1/oauth/authorize/google/callback", s.GoogleAuthorizeCallback)

	f.Get("/api/v1/oauth/me", middleware.JWT(), s.Me)
	f.Post("/api/v1/oauth/signout", middleware.JWT(), s.Signout)

	f.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{"code": fiber.StatusNotFound, "message": "No resource found"})
	})

	return f.Listen(":8080")
}
