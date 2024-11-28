package config

import (
	"log"

	"github.com/caarlos0/env/v11"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          int    `env:"PORT" default:"8080"`
	Dsn           string `env:"DSN" default:"localhost:3306"`
	JwtSecret     string `env:"JWT_SECRET"`
	JwtExpires    string `env:"JWT_EXPIRES"`
	RefreshSecret string `env:"REFRESH_SECRET"`
	RefreshExpiry string `env:"REFRESH_EXPIRY"`

	// SMTP
	SmtpHost     string `env:"SMTP_HOST"`
	SmtpPort     int    `env:"SMTP_PORT"`
	SmtpUser     string `env:"SMTP_USER"`
	SmtpPassword string `env:"SMTP_PASSWORD"`
	SmtpFrom     string `env:"SMTP_FROM"`
}

func New() *Config {
	if loadErr := godotenv.Load(".env"); loadErr != nil {
		log.Println("[Env]: unable to load .env file %v", loadErr)
	}

	var cfg Config

	if parseErr := env.Parse(&cfg); parseErr != nil {
		log.Println("[Env]: failed to parse environment variables: %v", parseErr)
	}

	return &cfg
}
