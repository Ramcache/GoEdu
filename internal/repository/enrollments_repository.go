package repository

import (
	"GoEdu/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepository interface {
	EnrollStudent(ctx context.Context, studentID, courseID int64) error
	GetStudentsByCourse(ctx context.Context, courseID int64) ([]*models.Student, error)
	UnEnrollStudent(ctx context.Context, studentID, courseID int64) error
	GetCoursesByStudent(ctx context.Context, studentID int64) ([]*models.Course, error)
}

type enrollmentRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) EnrollStudent(ctx context.Context, studentID, courseID int64) error {
	query := `
        INSERT INTO enrollments (student_id, course_id)
        VALUES ($1, $2)
        ON CONFLICT (student_id, course_id) DO NOTHING;
    `

	_, err := r.db.Exec(ctx, query, studentID, courseID)
	return err
}

func (r *enrollmentRepository) GetStudentsByCourse(ctx context.Context, courseID int64) ([]*models.Student, error) {
	query := `
        SELECT s.id, s.name, s.email
        FROM students s
        JOIN enrollments e ON s.id = e.student_id
        WHERE e.course_id = $1;
    `

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}
	return students, nil
}

func (r *enrollmentRepository) UnEnrollStudent(ctx context.Context, studentID, courseID int64) error {
	query := `
        DELETE FROM enrollments
        WHERE student_id = $1 AND course_id = $2;
    `
	_, err := r.db.Exec(ctx, query, studentID, courseID)
	return err
}

func (r *enrollmentRepository) GetCoursesByStudent(ctx context.Context, studentID int64) ([]*models.Course, error) {
	query := `
        SELECT c.id, c.name, c.description
        FROM courses c
        JOIN enrollments e ON c.id = e.course_id
        WHERE e.student_id = $1;
    `

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}
	return courses, nil
}
