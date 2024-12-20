package tests

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"testing"
//)
//
//func TestGetCoursesByInstructor(t *testing.T) {
//	ctx := context.Background()
//
//	_, err := db.Exec(ctx, "TRUNCATE TABLE courses, instructors RESTART IDENTITY CASCADE")
//	require.NoError(t, err, "Не удалось очистить таблицы")
//
//	_, err = db.Exec(ctx, "INSERT INTO instructors (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Преподаватель 1", "instructor1@domain.com", "securepassword")
//	require.NoError(t, err, "Не удалось добавить преподавателя 1")
//
//	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
//	require.NoError(t, err, "Не удалось добавить курс 1")
//
//	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 2, "Курс 2", "Описание курса 2", 1)
//	require.NoError(t, err, "Не удалось добавить курс 2")
//
//	testCases := []struct {
//		Name          string
//		Request       *proto.InstructorIDRequest
//		ExpectedCount int
//		ShouldError   bool
//		ExpectedCode  codes.Code
//	}{
//		{
//			Name:          "Успешное получение курсов",
//			Request:       &proto.InstructorIDRequest{InstructorId: 1},
//			ExpectedCount: 2,
//			ShouldError:   false,
//		},
//		{
//			Name:          "Преподаватель не найден",
//			Request:       &proto.InstructorIDRequest{InstructorId: 99},
//			ExpectedCount: 0,
//			ShouldError:   false,
//		},
//		{
//			Name:          "Некорректный ID преподавателя",
//			Request:       &proto.InstructorIDRequest{InstructorId: 0},
//			ExpectedCount: 0,
//			ShouldError:   true,
//			ExpectedCode:  codes.InvalidArgument,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.Name, func(t *testing.T) {
//			resp, err := client.GetCoursesByInstructor(ctx, tc.Request)
//
//			if tc.ShouldError {
//				require.Error(t, err, "Ожидалась ошибка, но её не было")
//				st, ok := status.FromError(err)
//				require.True(t, ok, "Ошибка не является статусной")
//				assert.Equal(t, tc.ExpectedCode, st.Code(), "Некорректный код ошибки")
//				t.Logf("Полученный код ошибки: %v", st.Code())
//				return
//			}
//
//			require.NoError(t, err, "Ошибка вызова GetCoursesByInstructor")
//			assert.NotNil(t, resp, "Ответ не должен быть пустым")
//			assert.Len(t, resp.Courses, tc.ExpectedCount, "Некорректное количество курсов в ответе")
//
//			// Проверка данных
//			if tc.ExpectedCount > 0 {
//				for _, course := range resp.Courses {
//					assert.NotZero(t, course.Id, "ID курса должен быть задан")
//					assert.NotEmpty(t, course.Name, "Название курса не должно быть пустым")
//					assert.NotEmpty(t, course.Description, "Описание курса не должно быть пустым")
//				}
//			}
//		})
//	}
//}
//
