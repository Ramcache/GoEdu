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

func TestRegisterStudent(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу студентов")

	testCases := []struct {
		Name         string
		Request      *proto.RegisterStudentRequest
		Expected     *proto.Student
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешная регистрация студента",
			Request: &proto.RegisterStudentRequest{
				Name:     "Студент 1",
				Email:    "student1_unique@domain.com",
				Password: "securepassword",
			},
			Expected: &proto.Student{
				Name:  "Студент 1",
				Email: "student1_unique@domain.com",
			},
			ShouldError: false,
		},
		{
			Name: "Студент с существующим email",
			Request: &proto.RegisterStudentRequest{
				Name:     "Студент 1",
				Email:    "student1@domain.com",
				Password: "securepassword",
			},
			ShouldError:  true,
			ExpectedCode: codes.AlreadyExists,
		},
	}

	_, err = db.Exec(ctx, "INSERT INTO students (name, email, password) VALUES ($1, $2, $3)", "Студент 1", "student1@domain.com", "hashedpassword")
	require.NoError(t, err, "Не удалось добавить существующего студента")

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientStudent.RegisterStudent(ctx, tc.Request)

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
			assert.Equal(t, tc.Expected.Name, resp.Name, "Имя студента не совпадает")
			assert.Equal(t, tc.Expected.Email, resp.Email, "Email студента не совпадает")
		})
	}
}

func TestLoginStudent(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу студентов")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)",
		1, "Студент 1", "student1@domain.com", "$2a$10$EHolyOdJvZejTUCtwWrUCu/bzLh1QK7ivkJpVRkjlC/YUkjmUzr86")
	require.NoError(t, err, "Не удалось добавить студента")

	testCases := []struct {
		Name         string
		Request      *proto.LoginRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешный вход студента",
			Request: &proto.LoginRequest{
				Email:    "student1@domain.com",
				Password: "securepassword",
			},
			ShouldError: false,
		},
		{
			Name: "Неверный пароль",
			Request: &proto.LoginRequest{
				Email:    "student1@domain.com",
				Password: "wrongpassword",
			},
			ShouldError:  true,
			ExpectedCode: codes.Unauthenticated,
		},
		{
			Name: "Студент не найден",
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
			resp, err := clientStudent.LoginStudent(ctx, tc.Request)

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
			assert.Equal(t, "student1@domain.com", resp.Email, "Email студента не совпадает")
		})
	}
}

func TestGetStudentProfile(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу студентов")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	testCases := []struct {
		Name         string
		Request      *proto.StudentIDRequest
		Expected     *proto.Student
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное получение профиля студента",
			Request: &proto.StudentIDRequest{
				Id: 1,
			},
			Expected: &proto.Student{
				Id:    1,
				Name:  "Студент 1",
				Email: "student1@domain.com",
				Token: "",
			},
			ShouldError: false,
		},
		{
			Name: "Студент не найден",
			Request: &proto.StudentIDRequest{
				Id: 99,
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientStudent.GetStudentProfile(ctx, tc.Request)

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
			assert.Equal(t, tc.Expected.Id, resp.Id, "ID студента не совпадает")
			assert.Equal(t, tc.Expected.Name, resp.Name, "Имя студента не совпадает")
			assert.Equal(t, tc.Expected.Email, resp.Email, "Email студента не совпадает")
		})
	}
}

func TestUpdateStudentProfile(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицу студентов")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	testCases := []struct {
		Name         string
		Request      *proto.UpdateStudentRequest
		Expected     *proto.Student
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное обновление профиля",
			Request: &proto.UpdateStudentRequest{
				Id:    1,
				Name:  "Обновлённый Студент",
				Email: "updated@domain.com",
			},
			Expected: &proto.Student{
				Id:    1,
				Name:  "Обновлённый Студент",
				Email: "updated@domain.com",
				Token: "",
			},
			ShouldError: false,
		},
		{
			Name: "Студент не найден",
			Request: &proto.UpdateStudentRequest{
				Id:    99,
				Name:  "Обновлённый Студент",
				Email: "updated@domain.com",
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientStudent.UpdateStudentProfile(ctx, tc.Request)

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
			assert.Equal(t, tc.Expected.Id, resp.Id, "ID студента не совпадает")
			assert.Equal(t, tc.Expected.Name, resp.Name, "Имя студента не совпадает")
			assert.Equal(t, tc.Expected.Email, resp.Email, "Email студента не совпадает")
		})
	}
}
