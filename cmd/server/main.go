package main

import (
	"net/http"
	"os"

	"go-api/cmd/server/auth"
	"go-api/internal/database"
	"go-api/internal/logger"
	"go-api/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load .env file
	godotenv.Load()
	logger.Logger.Info().Msg("Environment loaded")

	// JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		auth.JwtSecret = []byte(jwtSecret)
		logger.Logger.Info().Msg("JWT secret configured")
	} else {
		logger.Logger.Warn().Msg(" JWT_SECRET not set!")
	}

	// Connect to database
	err := database.Connect()
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg(" DB connection error")
	}
	logger.Logger.Info().Msg("Database connected")

	// Register routes
	router := routes.RegisterRoutes()
	logger.Logger.Info().Msg(" Routes registered")

	// Get port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Logger.Info().Str("port", port).Msg("Server starting")

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
