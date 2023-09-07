package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Jamess-Lucass/ecommerce-identity-service/services"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Me(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*services.Claim)

	uri := fmt.Sprintf("%s/api/v1/users/%s", os.Getenv("USER_SERVICE_BASE_URI"), claims.Subject)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	req.Header.Add("Authorization", c.Get("Authorization"))
	req.Header.Add("Cookie", c.Get("Cookie"))

	res, err := tracingClient.Do(req.WithContext(c.Context()))
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode >= fiber.StatusBadRequest {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}
		return c.Status(fiber.StatusBadRequest).JSON(string(bytes))
	}

	var user services.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
