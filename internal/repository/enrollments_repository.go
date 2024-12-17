package repository

import (
	"GoEdu/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepository interface {
	EnrollStudent(ctx context.Context, studentID, courseID int64) error
	GetStudentsByCourse(ctx context.Context, courseID int64) ([]*models.Student, error)
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
