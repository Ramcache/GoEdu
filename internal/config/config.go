package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	GRPCPort         string
	JWTSecretKey     string
	TokenExpiryHours int
}

type ConfigLoader interface {
	Load() (*Config, error)
}

type EnvConfigLoader struct{}

func (e *EnvConfigLoader) Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, использую переменные окружения")
	}

	expiryHours, err := strconv.Atoi(getEnv("TOKEN_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации TOKEN_EXPIRATION_HOURS: %v", err)
	}

	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/goedu"),
		GRPCPort:         getEnv("GRPC_PORT", "50051"),
		JWTSecretKey:     getEnv("JWT_SECRET_KEY", "your_jwt_secret_key"),
		TokenExpiryHours: expiryHours,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewConfig(loader ConfigLoader) *Config {
	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	return cfg
}
