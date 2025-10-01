package config

import "os"

type AppConfig struct {
	Port     string
	Env      string
	LogLevel string
}

func Load() AppConfig {
	return AppConfig{
		Port:     getEnv("APP_PORT", "5000"),
		Env:      getEnv("APP_ENV", "development"),
		LogLevel: getEnv("APP_LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
