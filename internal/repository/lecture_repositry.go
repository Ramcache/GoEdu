package repository

import (
	"GoEdu/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LectureRepository interface {
	AddLectureToCourse(ctx context.Context, lecture *models.Lecture) (*models.Lecture, error)
}

type lectureRepository struct {
	db *pgxpool.Pool
}

func NewLectureRepository(db *pgxpool.Pool) LectureRepository {
	return &lectureRepository{db: db}
}

func (r *lectureRepository) AddLectureToCourse(ctx context.Context, lecture *models.Lecture) (*models.Lecture, error) {
	query := `
        INSERT INTO lectures (course_id, title, content)
        VALUES ($1, $2, $3)
        RETURNING id, course_id, title, content;
    `

	var newLecture models.Lecture
	err := r.db.QueryRow(ctx, query, lecture.CourseID, lecture.Title, lecture.Content).
		Scan(&newLecture.ID, &newLecture.CourseID, &newLecture.Title, &newLecture.Content)
	if err != nil {
		return nil, err
	}

	return &newLecture, nil
}
