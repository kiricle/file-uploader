package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB_URL     string
	JWT_SECRET string
}

func SetupConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	DB_URL := os.Getenv("DATABASE_URL")
	JWT_SECRET := os.Getenv("JWT_SECRET")

	return &AppConfig{
		DB_URL:     DB_URL,
		JWT_SECRET: JWT_SECRET,
	}
}
