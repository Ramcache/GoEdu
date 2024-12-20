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

func TestEnrollStudent(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE enrollments, courses, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	testCases := []struct {
		Name         string
		Request      *proto.EnrollmentRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:        "Успешная запись студента на курс",
			Request:     &proto.EnrollmentRequest{StudentId: 1, CourseId: 1},
			ShouldError: false,
		},
		{
			Name:         "Некорректный ID студента",
			Request:      &proto.EnrollmentRequest{StudentId: 0, CourseId: 1},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Некорректный ID курса",
			Request:      &proto.EnrollmentRequest{StudentId: 1, CourseId: 0},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Запись на несуществующий курс",
			Request:      &proto.EnrollmentRequest{StudentId: 1, CourseId: 99},
			ShouldError:  true,
			ExpectedCode: codes.Internal,
		},
		{
			Name:         "Запись несуществующего студента на курс",
			Request:      &proto.EnrollmentRequest{StudentId: 99, CourseId: 1},
			ShouldError:  true,
			ExpectedCode: codes.Internal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEnrollments.EnrollStudent(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова EnrollStudent")
			assert.NotNil(t, resp, "Ответ должен быть непустым")

			// Проверка данных в базе
			var count int
			err = db.QueryRow(ctx, "SELECT COUNT(*) FROM enrollments WHERE student_id = $1 AND course_id = $2", tc.Request.StudentId, tc.Request.CourseId).Scan(&count)
			require.NoError(t, err, "Ошибка проверки данных в базе")
			assert.Equal(t, 1, count, "Ожидалась 1 запись в базе, но найдено %d")
		})
	}
}

func TestGetStudentsByCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE enrollments, courses, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента 1")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 2, "Студент 2", "student2@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента 2")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2), ($3, $4)", 1, 1, 2, 1)
	require.NoError(t, err, "Не удалось добавить записи о зачислении")

	testCases := []struct {
		Name             string
		Request          *proto.CourseIDRequest
		ExpectedCount    int
		ExpectedStudents []struct {
			ID    int64
			Name  string
			Email string
		}
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:          "Успешное получение студентов",
			Request:       &proto.CourseIDRequest{CourseId: 1},
			ExpectedCount: 2,
			ExpectedStudents: []struct {
				ID    int64
				Name  string
				Email string
			}{
				{1, "Студент 1", "student1@domain.com"},
				{2, "Студент 2", "student2@domain.com"},
			},
			ShouldError: false,
		},
		{
			Name:          "Курс не найден",
			Request:       &proto.CourseIDRequest{CourseId: 99},
			ExpectedCount: 0,
			ShouldError:   false,
		},
		{
			Name:         "Некорректный ID курса",
			Request:      &proto.CourseIDRequest{CourseId: 0},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEnrollments.GetStudentsByCourse(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова GetStudentsByCourse")
			assert.NotNil(t, resp, "Ответ должен быть непустым")
			assert.Len(t, resp.Students, tc.ExpectedCount, "Некорректное количество студентов")

			for i, student := range resp.Students {
				assert.Equal(t, tc.ExpectedStudents[i].ID, student.Id, "ID студента не совпадает")
				assert.Equal(t, tc.ExpectedStudents[i].Name, student.Name, "Имя студента не совпадает")
				assert.Equal(t, tc.ExpectedStudents[i].Email, student.Email, "Email студента не совпадает")
			}
		})
	}
}

func TestUnEnrollStudent(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE enrollments, courses, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2)", 1, 1)
	require.NoError(t, err, "Не удалось добавить запись о зачислении")

	testCases := []struct {
		Name         string
		Request      *proto.EnrollmentRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:        "Успешное удаление студента с курса",
			Request:     &proto.EnrollmentRequest{StudentId: 1, CourseId: 1},
			ShouldError: false,
		},
		{
			Name:         "Некорректный ID студента",
			Request:      &proto.EnrollmentRequest{StudentId: 0, CourseId: 1},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Некорректный ID курса",
			Request:      &proto.EnrollmentRequest{StudentId: 1, CourseId: 0},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:        "Попытка удалить несуществующего студента с курса",
			Request:     &proto.EnrollmentRequest{StudentId: 99, CourseId: 1},
			ShouldError: false,
		},
		{
			Name:        "Попытка удалить студента с несуществующего курса",
			Request:     &proto.EnrollmentRequest{StudentId: 1, CourseId: 99},
			ShouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEnrollments.UnEnrollStudent(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова UnEnrollStudent")
			assert.NotNil(t, resp, "Ответ должен быть непустым")

			var count int
			err = db.QueryRow(ctx, "SELECT COUNT(*) FROM enrollments WHERE student_id = $1 AND course_id = $2", tc.Request.StudentId, tc.Request.CourseId).Scan(&count)
			require.NoError(t, err, "Ошибка проверки данных в базе")
			assert.Equal(t, 0, count, "Ожидалось, что запись будет удалена из базы, но она найдена")
		})
	}
}

func TestGetCoursesByStudent(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE enrollments, courses, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс 1")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 2, "Курс 2", "Описание курса 2", 1)
	require.NoError(t, err, "Не удалось добавить курс 2")

	_, err = db.Exec(ctx, "INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2), ($1, $3)", 1, 1, 2)
	require.NoError(t, err, "Не удалось добавить записи о зачислении")

	testCases := []struct {
		Name            string
		Request         *proto.StudentIDRequest
		ExpectedCount   int
		ExpectedCourses []struct {
			ID          int64
			Name        string
			Description string
		}
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:          "Успешное получение курсов",
			Request:       &proto.StudentIDRequest{Id: 1},
			ExpectedCount: 2,
			ExpectedCourses: []struct {
				ID          int64
				Name        string
				Description string
			}{
				{1, "Курс 1", "Описание курса 1"},
				{2, "Курс 2", "Описание курса 2"},
			},
			ShouldError: false,
		},
		{
			Name:          "Студент не найден",
			Request:       &proto.StudentIDRequest{Id: 99},
			ExpectedCount: 0,
			ShouldError:   false,
		},
		{
			Name:         "Некорректный ID студента",
			Request:      &proto.StudentIDRequest{Id: 0},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEnrollments.GetCoursesByStudent(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова GetCoursesByStudent")
			assert.NotNil(t, resp, "Ответ должен быть непустым")
			assert.Len(t, resp.Courses, tc.ExpectedCount, "Некорректное количество курсов")

			for i, course := range resp.Courses {
				assert.Equal(t, tc.ExpectedCourses[i].ID, course.Id, "ID курса не совпадает")
				assert.Equal(t, tc.ExpectedCourses[i].Name, course.Name, "Название курса не совпадает")
				assert.Equal(t, tc.ExpectedCourses[i].Description, course.Description, "Описание курса не совпадает")
			}
		})
	}
}
