package tests

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"GoEdu/internal/repository"
	"GoEdu/internal/service"
	"GoEdu/proto"
)

var client proto.EducationServiceClient
var server *grpc.Server
var db *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Не удалось загрузить файл .env: %v", err)
		log.Println("Использую переменные окружения")
	}
	var err error
	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	courseRepo := repository.NewCourseRepository(db)
	educationService := service.NewEducationService(courseRepo)

	server = grpc.NewServer()
	proto.RegisterEducationServiceServer(server, educationService)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Ошибка сервера: %v", err)
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
}

func TestCreateCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Не удалось очистить таблицы: %v", err)
	}

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Тестовый преподаватель", "test@instructor.com", "securepassword")
	if err != nil {
		t.Fatalf("Не удалось добавить преподавателя: %v", err)
	}

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
		},
		{
			Name:         "Дублирующееся название курса",
			Request:      &proto.NewCourseRequest{Name: "Тестовый курс", Description: "Дублирующее название курса", InstructorId: 1},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
		},
		{
			Name:         "Пустое название курса",
			Request:      &proto.NewCourseRequest{Name: "", Description: "Курс без названия", InstructorId: 1},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
		},
		{
			Name:         "Пустое описание",
			Request:      &proto.NewCourseRequest{Name: "Курс без описания", Description: "", InstructorId: 1},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := client.CreateCourse(ctx, tc.Request)
			if tc.ShouldError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но ее не было")
				}
				return
			}

			if err != nil {
				t.Fatalf("Ошибка вызова CreateCourse: %v", err)
			}

			if resp.Id != tc.ExpectedID {
				t.Errorf("Ожидался ID %d, но получен %d", tc.ExpectedID, resp.Id)
			}

			if resp.Name != tc.ExpectedName {
				t.Errorf("Ожидалось название курса '%s', но получено '%s'", tc.ExpectedName, resp.Name)
			}

			var count int
			err = db.QueryRow(ctx, "SELECT COUNT(*) FROM courses WHERE name = $1", tc.Request.Name).Scan(&count)
			if err != nil {
				t.Fatalf("Ошибка проверки данных в базе: %v", err)
			}

			if !tc.ShouldError && count != 1 {
				t.Errorf("Ожидался 1 курс в базе, но найдено %d", count)
			}
		})
	}
}
