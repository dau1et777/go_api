package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port string
    JWT  string
}

func LoadConfig() *Config {
    _ = godotenv.Load()

    appPort := os.Getenv("APP_PORT")
    if appPort == "" {
        appPort = "8080"
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Println("⚠ WARNING: JWT_SECRET is empty!")
    }

    return &Config{
        Port: appPort,
        JWT:  jwtSecret,
    }
}
