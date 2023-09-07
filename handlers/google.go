package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Jamess-Lucass/ecommerce-identity-service/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthGoogleConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "openid", "profile"},
	Endpoint:     google.Endpoint,
}

type Response[T any] struct {
	Value []T `json:"value"`
}

type GoogleUser struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	Picture   string `json:"picture"`
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarUrl string `json:"avatarUrl"`
}

// GET: /api/v1/oauth/authorize/google
func (s *Server) RedirectGoogleAuthorize(c *fiber.Ctx) error {
	uri, err := url.Parse(oauthGoogleConfig.Endpoint.AuthURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Could not parse url",
		})
	}

	params := url.Values{}
	params.Add("client_id", oauthGoogleConfig.ClientID)
	params.Add("scope", strings.Join(oauthGoogleConfig.Scopes, " "))
	params.Add("redirect_uri", oauthGoogleConfig.RedirectURL)
	params.Add("response_type", "code")
	params.Add("state", c.Query("redirect_uri"))
	uri.RawQuery = params.Encode()

	return c.Status(fiber.StatusTemporaryRedirect).Redirect(uri.String())
}

// GET: /api/v1/oauth/authorize/google/callback
func (s *Server) GoogleAuthorizeCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := oauthGoogleConfig.Exchange(c.Context(), code)
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}

	s.logger.Sugar().Infof("google user :%v", googleUser)

	serviceJwt, _ := s.jwtService.CreateToken(services.User{Role: "Administrator"})

	uri, err := url.Parse(fmt.Sprintf("%s/api/v1/users", os.Getenv("USER_SERVICE_BASE_URI")))
	if err != nil {
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("email eq '%s'", googleUser.Email))
	uri.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", serviceJwt))
	res, err := tracingClient.Do(req.WithContext(c.Context()))
	if err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}
	defer res.Body.Close()

	var users Response[services.User]
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		s.logger.Sugar().Errorf("%v", err)
		return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
	}

	s.logger.Sugar().Infof("leng: %v", len(users.Value))

	if len(users.Value) == 0 {
		// Create user
		uri, err := url.Parse(fmt.Sprintf("%s/api/v1/users", os.Getenv("USER_SERVICE_BASE_URI")))
		if err != nil {
			return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
		}

		request := CreateUserRequest{
			Email:     googleUser.Email,
			FirstName: googleUser.FirstName,
			LastName:  "",
			AvatarUrl: googleUser.Picture,
		}

		body, _ := json.Marshal(request)
		req, err := http.NewRequest("POST", uri.String(), bytes.NewReader(body))
		if err != nil {
			return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", serviceJwt))

		response, err := tracingClient.Do(req.WithContext(c.Context()))
		if err != nil {
			s.logger.Sugar().Errorf("%v", err)
			return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
		}
		defer response.Body.Close()

		if response.StatusCode >= fiber.StatusBadRequest {
			bytes, _ := io.ReadAll(response.Body)
			s.logger.Sugar().Errorf("%v", string(bytes))
			return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
		}

		var user services.User
		if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
			s.logger.Sugar().Errorf("%v", err)
			return c.Redirect(fmt.Sprintf("%s?error=server_error", os.Getenv("LOGIN_UI_BASE_URL")))
		}

		jwt, _ := s.jwtService.CreateToken(user)

		cookie := new(fiber.Cookie)
		cookie.Name = "x-access-token"
		cookie.Value = jwt
		cookie.HTTPOnly = true
		cookie.Secure = true
		cookie.SameSite = "None"
		cookie.Expires = time.Now().Add(24 * time.Hour)

		if strings.EqualFold(os.Getenv("ENVIRONMENT"), "production") {
			cookie.Domain = "jameslucas.uk"
		} else {
			cookie.Domain = "localhost"
		}

		c.Cookie(cookie)

		return c.Status(fiber.StatusPermanentRedirect).Redirect(c.Query("state"))
	}

	jwt, _ := s.jwtService.CreateToken(users.Value[0])

	cookie := new(fiber.Cookie)
	cookie.Name = "x-access-token"
	cookie.Value = jwt
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = "None"
	cookie.Expires = time.Now().Add(24 * time.Hour)

	if strings.EqualFold(os.Getenv("ENVIRONMENT"), "production") {
		cookie.Domain = "jameslucas.uk"
	} else {
		cookie.Domain = "localhost"
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusPermanentRedirect).Redirect(c.Query("state"))
}
