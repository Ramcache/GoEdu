package main

import (
	"log"
	"net"

	"GoEdu/internal/config"
	"GoEdu/internal/repository"
	"GoEdu/internal/service"
	"GoEdu/proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	dbConfig := repository.Config{
		ConnectionString: cfg.DatabaseURL,
	}

	dbpool, err := repository.NewDBPool(dbConfig)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	courseRepo := repository.NewCourseRepository(dbpool)
	educationService := service.NewEducationService(courseRepo)

	server := grpc.NewServer()
	proto.RegisterEducationServiceServer(server, educationService)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	log.Printf("gRPC сервер запущен на :%s\n", cfg.GRPCPort)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
