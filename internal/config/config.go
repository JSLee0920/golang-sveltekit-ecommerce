package config

import "os"

type Config struct {
	AppEnv      string
	AppPort     string
	DatabaseURL string
	RedisURL    string
}

func Load() *Config {
	return &Config{
		AppEnv:      getEnv("APP_ENV"),
		AppPort:     getEnv("APP_PORT"),
		DatabaseURL: getEnv("DATABASE_URL"),
		RedisURL:    getEnv("REDIS_URL"),
	}
}

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return "Env not loading"
}
