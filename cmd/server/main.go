package main

import (
	"GoEdu/internal/config"
	"GoEdu/internal/middleware"
	"GoEdu/internal/repository"
	"GoEdu/internal/service"
	"GoEdu/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	loader := &config.EnvConfigLoader{}
	cfg := config.NewConfig(loader)

	dbConfig := repository.Config{
		ConnectionString: cfg.DatabaseURL,
	}
	dbpool, err := repository.NewDBPool(dbConfig)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbpool.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor([]byte(cfg.JWTSecretKey))),
	)

	// Репозитории
	enrollmentRepo := repository.NewEnrollmentRepository(dbpool)
	courseRepo := repository.NewCourseRepository(dbpool)
	studentRepo := repository.NewStudentRepository(dbpool)
	lectureRepo := repository.NewLectureRepository(dbpool)
	instructorRepo := repository.NewInstructorRepository(dbpool)
	reviewRepo := repository.NewReviewRepository(dbpool)

	// Сервисы
	enrollmentService := service.NewEnrollmentService(enrollmentRepo)
	educationService := service.NewEducationService(courseRepo)
	studentService := service.NewStudentService(studentRepo, cfg)
	lectureService := service.NewLectureService(lectureRepo)
	instructorService := service.NewInstructorService(instructorRepo, cfg)
	reviewService := service.NewReviewService(reviewRepo)

	// Регистрация сервисов
	proto.RegisterEducationServiceServer(grpcServer, educationService)
	proto.RegisterStudentServiceServer(grpcServer, studentService)
	proto.RegisterEnrollmentServiceServer(grpcServer, enrollmentService)
	proto.RegisterLectureServiceServer(grpcServer, lectureService)
	proto.RegisterInstructorServiceServer(grpcServer, instructorService)
	proto.RegisterReviewServiceServer(grpcServer, reviewService)

	// Запуск сервера
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	log.Printf("gRPC сервер запущен на :%s\n", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
