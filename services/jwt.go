package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTService struct {
}

type User struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatarUrl"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
}

type Claim struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) CreateToken(user User) (string, error) {
	claims := Claim{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
