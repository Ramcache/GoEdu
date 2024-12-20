package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"

	"google.golang.org/grpc/status"

	"GoEdu/proto"
)

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
			resp, err := clientEducation.CreateCourse(ctx, tc.Request)
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

func TestGetCourses(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING", 1, "Преподаватель 1", "instructor1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя 1")
	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING", 2, "Преподаватель 2", "instructor2@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя 2")

	_, err = db.Exec(ctx, "INSERT INTO courses (name, description, instructor_id) VALUES ($1, $2, $3)", "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс 1")
	_, err = db.Exec(ctx, "INSERT INTO courses (name, description, instructor_id) VALUES ($1, $2, $3)", "Курс 2", "Описание курса 2", 2)
	require.NoError(t, err, "Не удалось добавить курс 2")

	testCases := []struct {
		Name            string
		Request         *proto.Empty
		ExpectedCount   int
		ExpectedCourses []struct {
			Name        string
			Description string
		}
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:          "Успешное получение курсов",
			Request:       &proto.Empty{},
			ExpectedCount: 2,
			ExpectedCourses: []struct {
				Name        string
				Description string
			}{
				{"Курс 1", "Описание курса 1"},
				{"Курс 2", "Описание курса 2"},
			},
			ShouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEducation.GetCourses(ctx, tc.Request)
			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка при вызове GetCourses")

			assert.Len(t, resp.Courses, tc.ExpectedCount, "Некорректное количество курсов")

			for i, course := range resp.Courses {
				assert.Equal(t, tc.ExpectedCourses[i].Name, course.Name, "Название курса не совпадает")
				assert.Equal(t, tc.ExpectedCourses[i].Description, course.Description, "Описание курса не совпадает")
			}
		})
	}

	var count int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM courses").Scan(&count)
	require.NoError(t, err, "Ошибка при проверке количества курсов в базе")
	assert.Equal(t, 2, count, "Ожидалось, что в базе данных будет 2 курса, но найдено %d", count)
}

func TestGetCourseByID(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Преподаватель", "instructor@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	testCases := []struct {
		Name         string
		Request      *proto.CourseIDRequest
		ExpectedID   int64
		ExpectedName string
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:         "Успешное получение курса",
			Request:      &proto.CourseIDRequest{CourseId: 1},
			ExpectedID:   1,
			ExpectedName: "Курс 1",
			ShouldError:  false,
		},
		{
			Name:         "Курс не найден",
			Request:      &proto.CourseIDRequest{CourseId: 99},
			ExpectedID:   0,
			ExpectedName: "",
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEducation.GetCourseByID(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова GetCourseByID")
			assert.Equal(t, tc.ExpectedID, resp.Id, "ID курса не совпадает")
			assert.Equal(t, tc.ExpectedName, resp.Name, "Название курса не совпадает")

			var count int
			err = db.QueryRow(ctx, "SELECT COUNT(*) FROM courses WHERE id = $1", tc.Request.CourseId).Scan(&count)
			require.NoError(t, err, "Ошибка проверки данных в базе")
			assert.Equal(t, 1, count, "Ожидался 1 курс в базе, но найдено %d")
		})
	}
}

func TestUpdateCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Преподаватель", "instructor@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 2, "Курс 2", "Описание курса 2", 1)
	require.NoError(t, err, "Не удалось добавить второй курс")

	testCases := []struct {
		Name         string
		Request      *proto.UpdateCourseRequest
		ExpectedID   int64
		ExpectedName string
		ExpectedDesc string
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:         "Успешное обновление курса",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: "Обновленный курс", Description: "Обновленное описание курса"},
			ExpectedID:   1,
			ExpectedName: "Обновленный курс",
			ExpectedDesc: "Обновленное описание курса",
			ShouldError:  false,
		},
		{
			Name:         "Курс не найден",
			Request:      &proto.UpdateCourseRequest{Id: 99, Name: "Несуществующий курс", Description: "Описание"},
			ExpectedID:   0,
			ExpectedName: "",
			ExpectedDesc: "",
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
		{
			Name:         "Обновление курса без имени",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: "", Description: "Обновленное описание"},
			ExpectedID:   0,
			ExpectedName: "",
			ExpectedDesc: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Обновление курса без описания",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: "Обновленный курс", Description: ""},
			ExpectedID:   0,
			ExpectedName: "",
			ExpectedDesc: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Обновление с дублирующимся именем",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: "Курс 2", Description: "Дублирующее имя"},
			ExpectedID:   0,
			ExpectedName: "",
			ExpectedDesc: "",
			ShouldError:  true,
			ExpectedCode: codes.AlreadyExists,
		},
		{
			Name:         "Обновление без изменения данных",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: "Курс 1", Description: "Описание курса 1"},
			ExpectedID:   1,
			ExpectedName: "Курс 1",
			ExpectedDesc: "Описание курса 1",
			ShouldError:  false,
		},
		{
			Name:         "Обновление курса с длинным именем",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: string(make([]byte, 1001)), Description: "Обновленное описание"},
			ExpectedID:   0,
			ExpectedName: "",
			ExpectedDesc: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Обновление курса со специальными символами в имени",
			Request:      &proto.UpdateCourseRequest{Id: 1, Name: "Курс @$%!#*&", Description: "Описание курса"},
			ExpectedID:   1,
			ExpectedName: "Курс @$%!#*&",
			ExpectedDesc: "Описание курса",
			ShouldError:  false,
		},
		{
			Name:         "Обновление пустым запросом",
			Request:      &proto.UpdateCourseRequest{Id: 0, Name: "", Description: ""},
			ExpectedID:   0,
			ExpectedName: "",
			ExpectedDesc: "",
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEducation.UpdateCourse(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова UpdateCourse")
			assert.Equal(t, tc.ExpectedID, resp.Id, "ID курса не совпадает")
			assert.Equal(t, tc.ExpectedName, resp.Name, "Название курса не совпадает")
			assert.Equal(t, tc.ExpectedDesc, resp.Description, "Описание курса не совпадает")

			var name, description string
			err = db.QueryRow(ctx, "SELECT name, description FROM courses WHERE id = $1", tc.Request.Id).Scan(&name, &description)
			require.NoError(t, err, "Ошибка проверки данных в базе")
			assert.Equal(t, tc.ExpectedName, name, "Название курса в базе не совпадает")
			assert.Equal(t, tc.ExpectedDesc, description, "Описание курса в базе не совпадает")
		})
	}
}

func TestDeleteCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Преподаватель", "instructor@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	testCases := []struct {
		Name         string
		Request      *proto.CourseIDRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name:        "Успешное удаление курса",
			Request:     &proto.CourseIDRequest{CourseId: 1},
			ShouldError: false,
		},
		{
			Name:         "Курс не найден",
			Request:      &proto.CourseIDRequest{CourseId: 99},
			ShouldError:  true,
			ExpectedCode: codes.NotFound,
		},
		{
			Name:         "Ошибка в базе данных",
			Request:      &proto.CourseIDRequest{CourseId: -1},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name:         "Некорректный ID курса",
			Request:      &proto.CourseIDRequest{CourseId: -1},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEducation.DeleteCourse(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова DeleteCourse")
			assert.NotNil(t, resp, "Ответ должен быть непустым")

			var count int
			err = db.QueryRow(ctx, "SELECT COUNT(*) FROM courses WHERE id = $1", tc.Request.CourseId).Scan(&count)
			require.NoError(t, err, "Ошибка проверки данных в базе")
			assert.Equal(t, 0, count, "Курс не был удалён из базы")
		})
	}
}

func TestSearchCourses(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Преподаватель", "instructor@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить преподавателя")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Go Basics", "Основы языка Go", 1)
	require.NoError(t, err, "Не удалось добавить курс 1")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 2, "Advanced Go", "Продвинутый курс по Go", 1)
	require.NoError(t, err, "Не удалось добавить курс 2")

	testCases := []struct {
		Name          string
		Request       *proto.SearchRequest
		ExpectedCount int
		ExpectedNames []string
		ShouldError   bool
		ExpectedCode  codes.Code
	}{
		{
			Name:          "Поиск по ключевому слову 'Go'",
			Request:       &proto.SearchRequest{Keyword: "Go"},
			ExpectedCount: 2,
			ExpectedNames: []string{"Go Basics", "Advanced Go"},
			ShouldError:   false,
		},
		{
			Name:          "Поиск по ключевому слову 'Basics'",
			Request:       &proto.SearchRequest{Keyword: "Basics"},
			ExpectedCount: 1,
			ExpectedNames: []string{"Go Basics"},
			ShouldError:   false,
		},
		{
			Name:          "Поиск с пустым ключевым словом",
			Request:       &proto.SearchRequest{Keyword: ""},
			ExpectedCount: 0,
			ShouldError:   true,
			ExpectedCode:  codes.InvalidArgument,
		},
		{
			Name:          "Поиск без совпадений",
			Request:       &proto.SearchRequest{Keyword: "Python"},
			ExpectedCount: 0,
			ShouldError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientEducation.SearchCourses(ctx, tc.Request)

			if tc.ShouldError {
				require.Error(t, err, "Ожидалась ошибка, но её не было")
				st, ok := status.FromError(err)
				require.True(t, ok, "Ошибка не является статусной")
				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
				t.Logf("Полученный код ошибки: %v", st.Code())
				return
			}

			require.NoError(t, err, "Ошибка вызова SearchCourses")
			assert.NotNil(t, resp, "Ответ должен быть непустым")
			assert.Len(t, resp.Courses, tc.ExpectedCount, "Некорректное количество найденных курсов")

			for i, course := range resp.Courses {
				assert.Equal(t, tc.ExpectedNames[i], course.Name, "Название курса не совпадает")
			}
		})
	}
}
