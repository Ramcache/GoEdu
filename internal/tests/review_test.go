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

func TestAddReviewToCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE reviews, courses, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	testCases := []struct {
		Name         string
		Request      *proto.ReviewRequest
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное добавление отзыва",
			Request: &proto.ReviewRequest{
				StudentId: 1,
				CourseId:  1,
				Comment:   "Отличный курс!",
				Rating:    5,
			},
			ShouldError: false,
		},
		{
			Name: "Некорректный рейтинг",
			Request: &proto.ReviewRequest{
				StudentId: 1,
				CourseId:  1,
				Comment:   "Отзыв с некорректным рейтингом",
				Rating:    0,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
		{
			Name: "Пустой отзыв",
			Request: &proto.ReviewRequest{
				StudentId: 1,
				CourseId:  1,
				Comment:   "",
				Rating:    4,
			},
			ShouldError: false,
		},
		{
			Name: "Некорректный ID студента",
			Request: &proto.ReviewRequest{
				StudentId: 0,
				CourseId:  1,
				Comment:   "Некорректный студент",
				Rating:    4,
			},
			ShouldError:  true,
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := clientReview.AddReviewToCourse(ctx, tc.Request)

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

func TestGetReviewsByCourse(t *testing.T) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE TABLE reviews, courses, students RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Не удалось очистить таблицы")

	_, err = db.Exec(ctx, "INSERT INTO students (id, name, email, password) VALUES ($1, $2, $3, $4)", 1, "Студент 1", "student1@domain.com", "securepassword")
	require.NoError(t, err, "Не удалось добавить студента")

	_, err = db.Exec(ctx, "INSERT INTO courses (id, name, description, instructor_id) VALUES ($1, $2, $3, $4)", 1, "Курс 1", "Описание курса 1", 1)
	require.NoError(t, err, "Не удалось добавить курс")

	_, err = db.Exec(ctx, "INSERT INTO reviews (id, student_id, course_id, comment, rating, created_at) VALUES ($1, $2, $3, $4, $5, NOW())", 1, 1, 1, "Хороший курс", 4)
	require.NoError(t, err, "Не удалось добавить отзыв")

	testCases := []struct {
		Name         string
		Request      *proto.CourseIDRequest
		Expected     *proto.ReviewList
		ShouldError  bool
		ExpectedCode codes.Code
	}{
		{
			Name: "Успешное получение отзывов",
			Request: &proto.CourseIDRequest{
				CourseId: 1,
			},
			Expected: &proto.ReviewList{
				Reviews: []*proto.Review{
					{
						Id:        1,
						StudentId: 1,
						CourseId:  1,
						Comment:   "Хороший курс",
						Rating:    4,
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
			Expected: &proto.ReviewList{
				Reviews: []*proto.Review{},
			},
			ShouldError: false,
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
			resp, err := clientReview.GetReviewsByCourse(ctx, tc.Request)

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
			assert.Len(t, resp.Reviews, len(tc.Expected.Reviews), "Некорректное количество отзывов")

			for i, review := range resp.Reviews {
				assert.Equal(t, tc.Expected.Reviews[i].Id, review.Id, "ID отзыва не совпадает")
				assert.Equal(t, tc.Expected.Reviews[i].StudentId, review.StudentId, "ID студента не совпадает")
				assert.Equal(t, tc.Expected.Reviews[i].CourseId, review.CourseId, "ID курса не совпадает")
				assert.Equal(t, tc.Expected.Reviews[i].Comment, review.Comment, "Комментарий не совпадает")
				assert.Equal(t, tc.Expected.Reviews[i].Rating, review.Rating, "Рейтинг не совпадает")
			}
		})
	}
}
