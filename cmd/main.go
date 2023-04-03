package main

import (
	"github.com/Jamess-Lucass/ecommerce-identity-service/handlers"
	"github.com/Jamess-Lucass/ecommerce-identity-service/services"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Sugar().Warnf("could not flush: %v", err)
		}
	}()

	jwtService := services.NewJWTService()
	server := handlers.NewServer(logger, jwtService)

	if err := server.Start(); err != nil {
		logger.Sugar().Fatalf("error starting web server: %v", err)
	}
}
