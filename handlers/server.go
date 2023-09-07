package handlers

import (
	"net/http"

	"github.com/Jamess-Lucass/ecommerce-identity-service/services"
	"github.com/go-playground/validator/v10"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.uber.org/zap"
)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

type Server struct {
	validator  *validator.Validate
	logger     *zap.Logger
	jwtService *services.JWTService
}

func NewServer(logger *zap.Logger, jwtService *services.JWTService) *Server {
	return &Server{
		validator:  validator.New(),
		logger:     logger,
		jwtService: jwtService,
	}
}
