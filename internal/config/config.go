package config

import (
	"errors"
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	LogLevel string
	PostgresDSN string
	RedisAddr string
	RedisPass string
	CalcModelURL string
	
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	
	dsn := os.Getenv("POSTGRES_DSN")
	fmt.Println("dsn:", dsn)

	if dsn == "" {
		return nil, errors.New("POSTGRES_DSN environment variable is required")
	}
	
	
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "INFO"),
		
		PostgresDSN: dsn,
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
		CalcModelURL: getEnv("CALC_MODEL_URL", "http://localhost:8000/api/predict"),
	}, nil
}
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		
		return val
	}
	return fallback
	
}
