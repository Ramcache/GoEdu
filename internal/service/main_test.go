package service

import (
	"GoEdu/internal/config"
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"GoEdu/internal/logger"
	"GoEdu/internal/repository"
	"GoEdu/proto"
)

var (
	clientEducation   proto.EducationServiceClient
	clientEnrollments proto.EnrollmentServiceClient
	clientInstructor  proto.InstructorServiceClient
	clientLecture     proto.LectureServiceClient
	clientReview      proto.ReviewServiceClient
	clientStudent     proto.StudentServiceClient
	server            *grpc.Server
	db                *pgxpool.Pool
	zapLogger         *zap.Logger
)

func TestMain(m *testing.M) {
	var err error

	zapLogger, err = logger.NewLogger()
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			zapLogger.Error("Ошибка при завершении логгера", zap.Error(err))
		}
	}(zapLogger)
	zapLogger.Info("Инициализация тестов")

	if err := godotenv.Load("../../.env"); err != nil {
		zapLogger.Warn("Не удалось загрузить файл .env", zap.Error(err))
	}

	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		zapLogger.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
	}
	defer db.Close()
	zapLogger.Info("Успешное подключение к базе данных для тестов")

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		zapLogger.Fatal("Не удалось запустить сервер", zap.Error(err))
	}
	loader := &config.EnvConfigLoader{}
	cfg := config.NewConfig(loader)

	courseRepo := repository.NewCourseRepository(db)
	educationService := NewEducationService(db, courseRepo, zapLogger)

	enrollmentRepo := repository.NewEnrollmentRepository(db)
	enrollmentService := NewEnrollmentService(enrollmentRepo, zapLogger)

	instructorRepo := repository.NewInstructorRepository(db)
	instructorService := NewInstructorService(instructorRepo, cfg, zapLogger)

	lectureRepo := repository.NewLectureRepository(db)
	lectureService := NewLectureService(lectureRepo, zapLogger)

	reviewRepo := repository.NewReviewRepository(db)
	reviewService := NewReviewService(reviewRepo, zapLogger)

	studentRepo := repository.NewStudentRepository(db)
	studentService := NewStudentService(studentRepo, cfg, zapLogger)

	server = grpc.NewServer()
	proto.RegisterEducationServiceServer(server, educationService)
	proto.RegisterEnrollmentServiceServer(server, enrollmentService)
	proto.RegisterInstructorServiceServer(server, instructorService)
	proto.RegisterLectureServiceServer(server, lectureService)
	proto.RegisterReviewServiceServer(server, reviewService)
	proto.RegisterStudentServiceServer(server, studentService)

	go func() {
		if err := server.Serve(listener); err != nil {
			zapLogger.Fatal("Ошибка сервера", zap.Error(err))
		}
	}()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Ошибка при закрытии соединения: %v", err)
		}
	}(conn)

	clientEducation = proto.NewEducationServiceClient(conn)
	clientEnrollments = proto.NewEnrollmentServiceClient(conn)
	clientInstructor = proto.NewInstructorServiceClient(conn)
	clientLecture = proto.NewLectureServiceClient(conn)
	clientReview = proto.NewReviewServiceClient(conn)
	clientStudent = proto.NewStudentServiceClient(conn)

	code := m.Run()

	server.Stop()
	os.Exit(code)
}
