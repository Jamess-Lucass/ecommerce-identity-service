package main

import (
	"net/http"
	"os"

	"github.com/Jamess-Lucass/ecommerce-identity-service/handlers"
	"github.com/Jamess-Lucass/ecommerce-identity-service/services"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LOG_LEVEL = os.Getenv("LOG_LEVEL")
var LOG_LEVELS = map[string]zapcore.Level{
	"DEBUG": zap.DebugLevel,
	"INFO":  zap.InfoLevel,
	"WARN":  zap.WarnLevel,
	"ERROR": zap.ErrorLevel,
}

func main() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, LOG_LEVELS[LOG_LEVEL])
	logger := zap.New(core, zap.AddCaller())

	http.DefaultTransport = apmhttp.WrapRoundTripper(http.DefaultTransport, apmhttp.WithClientTrace())

	jwtService := services.NewJWTService()
	server := handlers.NewServer(logger, jwtService)

	if err := server.Start(); err != nil {
		logger.Sugar().Fatalf("error starting web server: %v", err)
	}
}
