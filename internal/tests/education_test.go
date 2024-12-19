package tests

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"log"
	"net"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"GoEdu/internal/logger"
	"GoEdu/internal/repository"
	"GoEdu/internal/service"
	"GoEdu/proto"
)

var client proto.EducationServiceClient
var server *grpc.Server
var db *pgxpool.Pool
var zapLogger *zap.Logger

func TestMain(m *testing.M) {
	var err error

	zapLogger, err = logger.NewLogger()
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	defer zapLogger.Sync()
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

	courseRepo := repository.NewCourseRepository(db)
	educationService := service.NewEducationService(db, courseRepo, zapLogger)

	server = grpc.NewServer()
	proto.RegisterEducationServiceServer(server, educationService)

	go func() {
		if err := server.Serve(listener); err != nil {
			zapLogger.Fatal("Ошибка сервера", zap.Error(err))
		}
	}()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Ошибка при закрытии соединения: %v", err)
		}
	}(conn)

	client = proto.NewEducationServiceClient(conn)

	code := m.Run()

	server.Stop()
	os.Exit(code)
}

type CreateCourseTestCase struct {
	Name         string
	Request      *proto.NewCourseRequest
	ExpectedID   int64
	ExpectedName string
	ShouldError  bool
	ExpectedCode codes.Code
}

func TestCreateCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Тестовый преподаватель", "test@instructor.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	require.NoError(t, err, "Не удалось начать транзакцию")
	defer tx.Rollback(ctx)

	testCases := []CreateCourseTestCase{
		{
			Name:         "Валидный курс",
			Request:      &proto.NewCourseRequest{Name: "Тестовый курс", Description: "Это тестовый курс", InstructorId: 1},
			ExpectedID:   1,
			ExpectedName: "Тестовый курс",
			ShouldError:  false,
		},
		{
			Name:         "Отсутствует преподаватель",
			Request:      &proto.NewCourseRequest{Name: "Тестовый курс 2", Description: "Это другой тестовый курс", InstructorId: 0},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Дублирующееся название курса",
			Request:      &proto.NewCourseRequest{Name: "Тестовый курс", Description: "Дублирующее название курса", InstructorId: 1},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
			ExpectedCode: codes.AlreadyExists,
		},

		{
			Name:         "Пустое название курса",
			Request:      &proto.NewCourseRequest{Name: "", Description: "Курс без названия", InstructorId: 1},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Пустое описание",
			Request:      &proto.NewCourseRequest{Name: "Курс без описания", Description: "", InstructorId: 1},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := client.CreateCourse(ctx, tc.Request)
			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова CreateCourse")
			assert.Equal(t, tc.ExpectedID, resp.Id, "ID курса не совпадает")
			assert.Equal(t, tc.ExpectedName, resp.Name, "Название курса не совпадает")

			var count int
			err = db.QueryRow(ctx, "SELECT COUNT(*) FROM courses WHERE name = $1", tc.Request.Name).Scan(&count)
			require.NoError(t, err, "Ошибка проверки данных в базе")
			assert.Equal(t, 1, count, "Ожидался 1 курс в базе, но найдено %d")
		})
	}
}
