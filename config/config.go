package config

import "os"

type Config struct {
	Port   string
	APIKey string
	DBPath string
}

func Load() *Config {
	return &Config{
		Port:   getEnv("PORT", "8080"),
		APIKey: getEnv("API_KEY", ""),
		DBPath: getEnv("DB_PATH", "banking.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
