package configs

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port     string
    MongoURI string
    JWTSecret string
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    return &Config{
        Port:      os.Getenv("PORT"),
        MongoURI:  os.Getenv("MONGO_URI"),
        JWTSecret: os.Getenv("JWT_SECRET"),
    }
}

