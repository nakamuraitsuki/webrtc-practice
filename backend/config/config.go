package config

import	"os"

type Config struct {
	Port string
	SecretKey string
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		SecretKey: getEnv("SECRET_KEY", "secret"),
	}
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}