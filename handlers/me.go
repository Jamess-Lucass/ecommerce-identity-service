package handlers

import (
	"encoding/json"
	"fmt"
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
		return c.Status(400).JSON(err.Error())
	}
	req.Header.Add("Authorization", c.Get("Authorization"))
	req.Header.Add("Cookie", c.Get("Cookie"))

	res, err := client.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Status(400).JSON(err.Error())
	}
	defer res.Body.Close()

	var user services.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(user)
}
