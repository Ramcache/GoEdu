package main

import (
	"GoEdu/internal/config"
	"GoEdu/internal/logger"
	"GoEdu/internal/middleware"
	"GoEdu/internal/repository"
	"GoEdu/internal/service"
	"GoEdu/proto"
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
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

	go func() {
		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithInsecure()}

		muxOptions := []runtime.ServeMuxOption{
			runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
				return metadata.Pairs("x-custom-header", req.Header.Get("X-Custom-Header"))
			}),
			runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
				http.Error(w, "Ошибка обработки: "+err.Error(), http.StatusInternalServerError)
			}),
		}

		gatewayMux := runtime.NewServeMux(muxOptions...)

		err := proto.RegisterEducationServiceHandlerFromEndpoint(ctx, gatewayMux, "localhost:"+cfg.GRPCPort, opts)
		if err != nil {
			zapLogger.Fatal("Не удалось зарегистрировать EducationService в gRPC Gateway", zap.Error(err))
		}

		router := mux.NewRouter()

		router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("R:/ProjectsGo/GoEdu/proto/swagger"))))
		router.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("R:/ProjectsGo/GoEdu/proto/swagger-ui"))))
		router.PathPrefix("/").Handler(gatewayMux)

		httpPort := cfg.HTTPPort
		zapLogger.Info("HTTP сервер запущен", zap.String("порт", httpPort))
		if err := http.ListenAndServe(":"+httpPort, router); err != nil {
			zapLogger.Fatal("Ошибка запуска HTTP сервера", zap.Error(err))
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor([]byte(cfg.JWTSecretKey), zapLogger)),
	)

	reflection.Register(grpcServer)

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

	// Запуск gRPC сервера
	go func() {
		listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			zapLogger.Fatal("Не удалось запустить сервер", zap.Error(err))
		}

		zapLogger.Info("gRPC сервер запущен", zap.String("порт", cfg.GRPCPort))
		if err := grpcServer.Serve(listener); err != nil {
			zapLogger.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	select {}
}
