package repository

import (
	"GoEdu/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LectureRepository interface {
	AddLectureToCourse(ctx context.Context, lecture *models.Lecture) (*models.Lecture, error)
	GetLecturesByCourse(ctx context.Context, courseID int64) ([]*models.Lecture, error)
	GetLectureContent(ctx context.Context, lectureID int64) (*models.Lecture, error)
	UpdateLecture(ctx context.Context, lecture *models.Lecture) (*models.Lecture, error)
	DeleteLecture(ctx context.Context, lectureID int64) error
	MarkLectureAsCompleted(ctx context.Context, studentID, lectureID int64) error
	GetCourseProgress(ctx context.Context, studentID, courseID int64) (int32, error)
	GetRecommendedCourses(ctx context.Context, studentID int64) ([]*models.Course, error)
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

func (r *lectureRepository) UpdateLecture(ctx context.Context, lecture *models.Lecture) (*models.Lecture, error) {
	query := `
        UPDATE lectures
        SET title = COALESCE(NULLIF($1, ''), title),
            content = COALESCE(NULLIF($2, ''), content)
        WHERE id = $3
        RETURNING id, course_id, title, content;
    `

	var updatedLecture models.Lecture
	err := r.db.QueryRow(ctx, query, lecture.Title, lecture.Content, lecture.ID).
		Scan(&updatedLecture.ID, &updatedLecture.CourseID, &updatedLecture.Title, &updatedLecture.Content)
	if err != nil {
		return nil, err
	}

	return &updatedLecture, nil
}

func (r *lectureRepository) DeleteLecture(ctx context.Context, lectureID int64) error {
	query := `
        DELETE FROM lectures
        WHERE id = $1;
    `

	commandTag, err := r.db.Exec(ctx, query, lectureID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *lectureRepository) MarkLectureAsCompleted(ctx context.Context, studentID, lectureID int64) error {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM lectures WHERE id = $1)`, lectureID).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("лекция с ID %d не найдена", lectureID)
	}

	err = r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)`, studentID).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("студент с ID %d не найден", studentID)
	}

	query := `
        INSERT INTO lecture_completions (student_id, lecture_id)
        VALUES ($1, $2)
        ON CONFLICT (student_id, lecture_id) DO NOTHING;
    `
	_, err = r.db.Exec(ctx, query, studentID, lectureID)
	return err
}

func (r *lectureRepository) GetCourseProgress(ctx context.Context, studentID, courseID int64) (int32, error) {
	var totalLectures int
	var completedLectures int

	err := r.db.QueryRow(ctx, `
        SELECT COUNT(*) FROM lectures WHERE course_id = $1;
    `, courseID).Scan(&totalLectures)
	if err != nil {
		return 0, err
	}

	if totalLectures == 0 {
		return 0, fmt.Errorf("в курсе нет лекций")
	}

	err = r.db.QueryRow(ctx, `
        SELECT COUNT(*) FROM lecture_completions lc
        JOIN lectures l ON lc.lecture_id = l.id
        WHERE lc.student_id = $1 AND l.course_id = $2;
    `, studentID, courseID).Scan(&completedLectures)
	if err != nil {
		return 0, err
	}

	progress := (float64(completedLectures) / float64(totalLectures)) * 100
	return int32(progress), nil
}

func (r *lectureRepository) GetRecommendedCourses(ctx context.Context, studentID int64) ([]*models.Course, error) {
	query := `
        SELECT DISTINCT c.id, c.name, c.description, c.instructor_id
        FROM courses c
        WHERE c.id NOT IN (
            SELECT e.course_id
            FROM enrollments e
            WHERE e.student_id = $1
        )
        ORDER BY RANDOM()
        LIMIT 5;
    `

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.InstructorID); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	return courses, nil
}
