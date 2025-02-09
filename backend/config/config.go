package config

import (
	"os"
	"time"
)

type Config struct {
	Port string
	SecretKey string
	TokenExpiry time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		SecretKey: getEnv("SECRET_KEY", "secret"),
		TokenExpiry: paraseDuration(getEnv("TOKEN_EXPIRY", "24h")),
	}
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func paraseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return d
}