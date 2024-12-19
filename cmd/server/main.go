package main

import (
	"GoEdu/internal/config"
	"GoEdu/internal/logger"
	"GoEdu/internal/middleware"
	"GoEdu/internal/repository"
	"GoEdu/internal/service"
	"GoEdu/proto"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	defer zapLogger.Sync()
	zapLogger.Info("Инициализация сервера")

	loader := &config.EnvConfigLoader{}
	cfg := config.NewConfig(loader)

	zapLogger.Info("Настройка подключения к базе данных")
	dbConfig := repository.Config{
		ConnectionString: cfg.DatabaseURL,
	}
	dbpool, err := repository.NewDBPool(dbConfig)
	if err != nil {
		zapLogger.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
	}
	defer dbpool.Close()
	zapLogger.Info("Успешное подключение к базе данных")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor([]byte(cfg.JWTSecretKey), zapLogger)),
	)

	// Репозитории
	enrollmentRepo := repository.NewEnrollmentRepository(dbpool)
	courseRepo := repository.NewCourseRepository(dbpool)
	studentRepo := repository.NewStudentRepository(dbpool)
	lectureRepo := repository.NewLectureRepository(dbpool)
	instructorRepo := repository.NewInstructorRepository(dbpool)
	reviewRepo := repository.NewReviewRepository(dbpool)

	// Сервисы
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, zapLogger)
	educationService := service.NewEducationService(dbpool, courseRepo, zapLogger)
	studentService := service.NewStudentService(studentRepo, cfg, zapLogger)
	lectureService := service.NewLectureService(lectureRepo, zapLogger)
	instructorService := service.NewInstructorService(instructorRepo, cfg, zapLogger)
	reviewService := service.NewReviewService(reviewRepo, zapLogger)

	// Регистрация сервисов
	proto.RegisterEducationServiceServer(grpcServer, educationService)
	proto.RegisterStudentServiceServer(grpcServer, studentService)
	proto.RegisterEnrollmentServiceServer(grpcServer, enrollmentService)
	proto.RegisterLectureServiceServer(grpcServer, lectureService)
	proto.RegisterInstructorServiceServer(grpcServer, instructorService)
	proto.RegisterReviewServiceServer(grpcServer, reviewService)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		zapLogger.Fatal("Не удалось запустить сервер", zap.Error(err))
	}

	zapLogger.Info("gRPC сервер запущен", zap.String("порт", cfg.GRPCPort))
	if err := grpcServer.Serve(listener); err != nil {
		zapLogger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
