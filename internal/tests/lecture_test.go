package tests

import (
	"GoEdu/proto"
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestAddLectureToCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE lectures, courses RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	testCases := []struct {
		Name         string
		Request      *proto.LectureRequest
		Expected     *proto.Lecture
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное добавление лекции",
			Request: &proto.LectureRequest{
				CourseId: 1,
				Title:    "Лекция 1",
				Content:  "Содержание лекции 1",
			},
			Expected: &proto.Lecture{
				CourseId: 1,
				Title:    "Лекция 1",
				Content:  "Содержание лекции 1",
			},
			ShouldError: false,
		},
		{
			Name: "Некорректный запрос: пустое название лекции",
			Request: &proto.LectureRequest{
				CourseId: 1,
				Title:    "",
				Content:  "Содержание лекции 1",
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.AddLectureToCourse(ctx, tc.Request)

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

			assert.Equal(t, tc.Expected.CourseId, resp.CourseId, "ID курса не совпадает")
			assert.Equal(t, tc.Expected.Title, resp.Title, "Название лекции не совпадает")
			assert.Equal(t, tc.Expected.Content, resp.Content, "Содержание лекции не совпадает")
		})
	}
}

func TestGetLecturesByCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE lectures, courses RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO lectures (id, course_id, title, content) VALUES ($1, $2, $3, $4)", 1, 1, "Лекция 1", "Содержание лекции 1")
	require.NoError(t, err, "Не удалось добавить лекцию")

	testCases := []struct {
		Name         string
		Request      *proto.CourseIDRequest
		Expected     *proto.LectureList
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное получение лекций",
			Request: &proto.CourseIDRequest{
				CourseId: 1,
			},
			Expected: &proto.LectureList{
				Lectures: []*proto.Lecture{
					{
						Id:       1,
						CourseId: 1,
						Title:    "Лекция 1",
						Content:  "Содержание лекции 1",
					},
				},
			},
			ShouldError: false,
		},
		{
			Name: "Курс не найден",
			Request: &proto.CourseIDRequest{
				CourseId: 99,
			},
			ShouldError: false,
			Expected:    &proto.LectureList{Lectures: []*proto.Lecture{}},
		},
		{
			Name: "Некорректный ID курса",
			Request: &proto.CourseIDRequest{
				CourseId: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.GetLecturesByCourse(ctx, tc.Request)

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
			assert.Len(t, resp.Lectures, len(tc.Expected.Lectures), "Некорректное количество лекций")

			for i, lecture := range resp.Lectures {
				assert.Equal(t, tc.Expected.Lectures[i].Id, lecture.Id, "ID лекции не совпадает")
				assert.Equal(t, tc.Expected.Lectures[i].CourseId, lecture.CourseId, "ID курса не совпадает")
				assert.Equal(t, tc.Expected.Lectures[i].Title, lecture.Title, "Название лекции не совпадает")
				assert.Equal(t, tc.Expected.Lectures[i].Content, lecture.Content, "Содержание лекции не совпадает")
			}
		})
	}
}

func TestGetLectureContent(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE lectures RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO lectures (id, course_id, title, content) VALUES ($1, $2, $3, $4)", 1, 1, "Лекция 1", "Содержание лекции 1")
	require.NoError(t, err, "Не удалось добавить лекцию")

	testCases := []struct {
		Name         string
		Request      *proto.LectureIDRequest
		Expected     *proto.LectureContent
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное получение содержания лекции",
			Request: &proto.LectureIDRequest{
				LectureId: 1,
			},
			Expected: &proto.LectureContent{
				Id:      1,
				Content: "Содержание лекции 1",
			},
			ShouldError: false,
		},
		{
			Name: "Лекция не найдена",
			Request: &proto.LectureIDRequest{
				LectureId: 99,
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
		{
			Name: "Некорректный ID лекции",
			Request: &proto.LectureIDRequest{
				LectureId: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.GetLectureContent(ctx, tc.Request)

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

			assert.Equal(t, tc.Expected.Id, resp.Id, "ID лекции не совпадает")
			assert.Equal(t, tc.Expected.Content, resp.Content, "Содержание лекции не совпадает")
		})
	}
}

func TestDeleteLecture(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE lectures RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO lectures (id, course_id, title, content) VALUES ($1, $2, $3, $4)", 1, 1, "Лекция 1", "Содержание лекции 1")
	require.NoError(t, err, "Не удалось добавить лекцию")

	testCases := []struct {
		Name         string
		Request      *proto.LectureIDRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное удаление лекции",
			Request: &proto.LectureIDRequest{
				LectureId: 1,
			},
			ShouldError: false,
		},
		{
			Name: "Лекция не найдена",
			Request: &proto.LectureIDRequest{
				LectureId: 99,
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
		{
			Name: "Некорректный ID лекции",
			Request: &proto.LectureIDRequest{
				LectureId: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.DeleteLecture(ctx, tc.Request)

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
		})
	}
}

func TestGetRecommendedCourses(t *testing.T) {
	ctx := context.Background()

	// Очистка таблиц
	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, enrollments, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	// Добавление преподавателей
	_, err = db.Exec(ctx, `
		INSERT INTO instructors (id, name, email, password) 
		VALUES 
		(1, 'Инструктор 1', 'instructor1@domain.com', 'password1'),
		(2, 'Инструктор 2', 'instructor2@domain.com', 'password2'),
		(3, 'Инструктор 3', 'instructor3@domain.com', 'password3')
	`)
	require.NoError(t, err, "Не удалось добавить преподавателей")

	// Добавление курсов
	_, err = db.Exec(ctx, `
		INSERT INTO courses (id, name, description, instructor_id) 
		VALUES 
		(1, 'Курс 1', 'Описание курса 1', 1),
		(2, 'Курс 2', 'Описание курса 2', 2),
		(3, 'Курс 3', 'Описание курса 3', 3)
	`)
	require.NoError(t, err, "Не удалось добавить курсы")

	// Добавление зачисления
	_, err = db.Exec(ctx, "INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2)", 1, 1)
	require.NoError(t, err, "Не удалось добавить зачисления")

	testCases := []struct {
		Name         string
		Request      *proto.StudentIDRequest
		Expected     *proto.CourseList
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное получение рекомендованных курсов",
			Request: &proto.StudentIDRequest{
				Id: 1,
			},
			Expected: &proto.CourseList{
				Courses: []*proto.Course{
					{
						Id:           2,
						Name:         "Курс 2",
						Description:  "Описание курса 2",
						InstructorId: 2,
					},
					{
						Id:           3,
						Name:         "Курс 3",
						Description:  "Описание курса 3",
						InstructorId: 3,
					},
				},
			},
			ShouldError: false,
		},
		{
			Name: "Некорректный ID студента",
			Request: &proto.StudentIDRequest{
				Id: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.GetRecommendedCourses(ctx, tc.Request)

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
			assert.Len(t, resp.Courses, len(tc.Expected.Courses), "Некорректное количество курсов")

			for i, course := range resp.Courses {
				assert.Equal(t, tc.Expected.Courses[i].Id, course.Id, "ID курса не совпадает")
				assert.Equal(t, tc.Expected.Courses[i].Name, course.Name, "Название курса не совпадает")
				assert.Equal(t, tc.Expected.Courses[i].Description, course.Description, "Описание курса не совпадает")
				assert.Equal(t, tc.Expected.Courses[i].InstructorId, course.InstructorId, "ID инструктора не совпадает")
			}
		})
	}
}

func TestGetCourseProgress(t *testing.T) {
	ctx := context.Background()

	// Подготовка базы данных
	_, err := db.Exec(ctx, "TRUNCATE TABLE lecture_completions, lectures, courses RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO lectures (id, course_id, title, content) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)", 1, 1, "Лекция 1", "Содержание лекции 1", 2, 1, "Лекция 2", "Содержание лекции 2")
	require.NoError(t, err, "Не удалось добавить лекции")

	_, err = db.Exec(ctx, "INSERT INTO lecture_completions (lecture_id, student_id) VALUES ($1, $2)", 1, 1)
	require.NoError(t, err, "Не удалось добавить завершенные лекции")

	testCases := []struct {
		Name         string
		Request      *proto.CourseProgressRequest
		Expected     *proto.CourseProgress
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное получение прогресса курса",
			Request: &proto.CourseProgressRequest{
				CourseId:  1,
				StudentId: 1,
			},
			Expected: &proto.CourseProgress{
				CourseId:         1,
				StudentId:        1,
				CompletedPercent: 50,
			},
			ShouldError: false,
		},
		{
			Name: "Курс не найден",
			Request: &proto.CourseProgressRequest{
				CourseId:  99,
				StudentId: 1,
			},
			ShouldError:  true,
			ExpectedCode: codes.Internal,
		},
		{
			Name: "Некорректный ID студента",
			Request: &proto.CourseProgressRequest{
				CourseId:  1,
				StudentId: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name: "Некорректный ID курса",
			Request: &proto.CourseProgressRequest{
				CourseId:  0,
				StudentId: 1,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.GetCourseProgress(ctx, tc.Request)

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

			assert.Equal(t, tc.Expected.CourseId, resp.CourseId, "ID курса не совпадает")
			assert.Equal(t, tc.Expected.StudentId, resp.StudentId, "ID студента не совпадает")
			assert.Equal(t, tc.Expected.CompletedPercent, resp.CompletedPercent, "Прогресс курса не совпадает")
		})
	}
}

func TestMarkLectureAsCompleted(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE lecture_completions, lectures, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	_, err = db.Exec(ctx, "INSERT INTO lectures (id, course_id, title, content) VALUES ($1, $2, $3, $4)", 1, 1, "Лекция 1", "Содержание лекции 1")
	require.NoError(t, err, "Не удалось добавить лекцию")

	testCases := []struct {
		Name         string
		Request      *proto.LectureCompletionRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешная отметка лекции как завершённой",
			Request: &proto.LectureCompletionRequest{
				StudentId: 1,
				LectureId: 1,
			},
			ShouldError: false,
		},
		{
			Name: "Лекция не найдена",
			Request: &proto.LectureCompletionRequest{
				StudentId: 1,
				LectureId: 99,
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
		{
			Name: "Студент не найден",
			Request: &proto.LectureCompletionRequest{
				StudentId: 99,
				LectureId: 1,
			},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
		{
			Name: "Некорректный запрос",
			Request: &proto.LectureCompletionRequest{
				StudentId: 0,
				LectureId: 0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientLecture.MarkLectureAsCompleted(ctx, tc.Request)

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
		})
	}
}
