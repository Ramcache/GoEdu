package service

import (
	"GoEdu/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestRegisterInstructor(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (name, email, password) VALUES ($1, $2, $3)", "Преподаватель 1", "instructor1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	testCases := []struct {
		Name         string
		Request      *proto.RegisterInstructorRequest
		Expected     *proto.Instructor
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешная регистрация преподавателя",
			Request: &proto.RegisterInstructorRequest{
				Name:     "Преподаватель 2",
				Email:    "instructor2@domain.com",
				Password: "securepassword",
			},
			Expected: &proto.Instructor{
				Name:  "Преподаватель 2",
				Email: "instructor2@domain.com",
			},
			ShouldError: false,
		},
		{
			Name: "Преподаватель с существующим email",
			Request: &proto.RegisterInstructorRequest{
				Name:     "Преподаватель 1",
				Email:    "instructor1@domain.com",
				Password: "securepassword",
			},
			ShouldError:  true,
			ExpectedCode: codes.AlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientInstructor.RegisterInstructor(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова метода")
			assert.NotNil(t, resp, "Ответ должен быть непустым")

			assert.Equal(t, tc.Expected.Name, resp.Name, "Имя преподавателя не совпадает")
			assert.Equal(t, tc.Expected.Email, resp.Email, "Email преподавателя не совпадает")
		})
	}
}

func TestLoginInstructor(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу преподавателей")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)",
		1, "Препод 1", "prepod1@domain.com", "$2a$10$EHolyOdJvZejTUCtwWrUCu/bzLh1QK7ivkJpVRkjlC/YUkjmUzr86")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	testCases := []struct {
		Name         string
		Request      *proto.LoginRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешный вход преподавателя",
			Request: &proto.LoginRequest{
				Email:    "prepod1@domain.com",
				Password: "securepassword",
			},
			ShouldError: false,
		},
		{
			Name: "Неверный пароль",
			Request: &proto.LoginRequest{
				Email:    "prepod1@domain.com",
				Password: "wrongpassword",
			},
			ShouldError:  true,
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "Преподаватель не найден",
			Request: &proto.LoginRequest{
				Email:    "unknown@domain.com",
				Password: "securepassword",
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientInstructor.LoginInstructor(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова метода")
			assert.NotNil(t, resp, "Ответ должен быть непустым")
			assert.Equal(t, "prepod1@domain.com", resp.Email, "Email преподавателя не совпадает")
			assert.Equal(t, "Препод 1", resp.Name, "Имя преподавателя не совпадает")
			assert.NotEmpty(t, resp.Token, "Токен не должен быть пустым")
		})

	}
}

func TestGetCoursesByInstructor(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE instructors, courses RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Преподаватель 1", "instructor1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс 1")

	testCases := []struct {
		Name         string
		Request      *proto.InstructorIDRequest
		Expected     *proto.CourseList
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Получение курсов преподавателя",
			Request: &proto.InstructorIDRequest{
				InstructorId: 1,
			},
			Expected: &proto.CourseList{
				Courses: []*proto.Course{
					{
						Id:           1,
						Name:         "Курс 1",
						Description:  "Описание курса 1",
						InstructorId: 1,
					},
				},
			},
			ShouldError: false,
		},
		{
			Name: "Некорректный ID преподавателя",
			Request: &proto.InstructorIDRequest{
				InstructorId: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientInstructor.GetCoursesByInstructor(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова метода")
			assert.NotNil(t, resp, "Ответ должен быть непустым")

			assert.Len(t, resp.Courses, len(tc.Expected.Courses), "Количество курсов не совпадает")
			for i, course := range resp.Courses {
				assert.Equal(t, tc.Expected.Courses[i].Id, course.Id, "ID курса не совпадает")
				assert.Equal(t, tc.Expected.Courses[i].Name, course.Name, "Название курса не совпадает")
				assert.Equal(t, tc.Expected.Courses[i].Description, course.Description, "Описание курса не совпадает")
				assert.Equal(t, tc.Expected.Courses[i].InstructorId, course.InstructorId, "ID преподавателя не совпадает")
			}
		})
	}
}
