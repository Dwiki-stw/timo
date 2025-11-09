package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment")
	}

	return &Config{
		App: App{
			Host: os.Getenv("APP_HOST"),
			Port: os.Getenv("APP_PORT"),
		},
		DB: DB{
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SslMode:  os.Getenv("DB_SSLMODE"),
		},
		JwtKey: []byte(os.Getenv("JWT_KEY")),
	}
}
