package config

import (
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

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, использую переменные окружения")
	}

	expiryHours, err := strconv.Atoi(getEnv("TOKEN_EXPIRATION_HOURS", "24"))
	if err != nil {
		log.Fatalf("Ошибка конвертации TOKEN_EXPIRATION_HOURS: %v", err)
	}

	return Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/goedu"),
		GRPCPort:         getEnv("GRPC_PORT", "50051"),
		JWTSecretKey:     getEnv("JWT_SECRET_KEY", "your_jwt_secret_key"),
		TokenExpiryHours: expiryHours,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
