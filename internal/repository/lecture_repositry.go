package repository

import (
	"GoEdu/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LectureRepository interface {
	AddLectureToCourse(ctx context.Context, lecture *models.Lecture) (*models.Lecture, error)
	GetLecturesByCourse(ctx context.Context, courseID int64) ([]*models.Lecture, error)
	GetLectureContent(ctx context.Context, lectureID int64) (*models.Lecture, error)
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

func (r *lectureRepository) GetLecturesByCourse(ctx context.Context, courseID int64) ([]*models.Lecture, error) {
	query := `
        SELECT id, course_id, title, content
        FROM lectures
        WHERE course_id = $1;
    `

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lectures []*models.Lecture
	for rows.Next() {
		var lecture models.Lecture
		if err := rows.Scan(&lecture.ID, &lecture.CourseID, &lecture.Title, &lecture.Content); err != nil {
			return nil, err
		}
		lectures = append(lectures, &lecture)
	}
	return lectures, nil
}

func (r *lectureRepository) GetLectureContent(ctx context.Context, lectureID int64) (*models.Lecture, error) {
	query := `
        SELECT id, course_id, title, content
        FROM lectures
        WHERE id = $1;
    `

	var lecture models.Lecture
	err := r.db.QueryRow(ctx, query, lectureID).Scan(
		&lecture.ID, &lecture.CourseID, &lecture.Title, &lecture.Content,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &lecture, nil
}
