package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	GRPCPort    string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, использую переменные окружения")
	}

	return Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		GRPCPort:    getEnv("GRPC_PORT", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
