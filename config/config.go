package config

import "os"

type Config struct {
	Port           string
	APIKey         string
	DBPath         string
	FrontendOrigin string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		APIKey:         getEnv("API_KEY", ""),
		DBPath:         getEnv("DB_PATH", "banking.db"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
