package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Config struct {
	ConnectionString string
}

func NewDBPool(cfg Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}

	log.Println("Успешно подключено к базе данных")
	return dbpool, nil
}
