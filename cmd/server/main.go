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
	instructorService := service.NewInstructorService(instructorRepo)
	reviewService := service.NewReviewService(reviewRepo)

	// gRPC сервер
	server := grpc.NewServer()
	proto.RegisterEducationServiceServer(server, educationService)
	proto.RegisterStudentServiceServer(server, studentService)
	proto.RegisterEnrollmentServiceServer(server, enrollmentService)
	proto.RegisterLectureServiceServer(server, lectureService)
	proto.RegisterInstructorServiceServer(server, instructorService)
	proto.RegisterReviewServiceServer(server, reviewService)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	log.Printf("gRPC сервер запущен на :%s\n", cfg.GRPCPort)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
