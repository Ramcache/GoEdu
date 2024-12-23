package repository

import (
	"GoEdu/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InstructorRepository interface {
	RegisterInstructor(ctx context.Context, instructor *models.Instructor) (int64, error)
	GetInstructorByEmail(ctx context.Context, email string) (*models.Instructor, error)
	GetCoursesByInstructor(ctx context.Context, instructorID int64) ([]*models.Course, error)
	GetInstructorByID(ctx context.Context, id int64) (*models.Instructor, error)
	UpdateInstructor(ctx context.Context, instructor *models.Instructor) error
}

type instructorRepository struct {
	db *pgxpool.Pool
}

func NewInstructorRepository(db *pgxpool.Pool) InstructorRepository {
	return &instructorRepository{db: db}
}

var ErrInstructorNotFound = errors.New("преподаватель не найден")

func (r *instructorRepository) RegisterInstructor(ctx context.Context, instructor *models.Instructor) (int64, error) {
	query := `
        INSERT INTO instructors (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	var id int64
	err := r.db.QueryRow(ctx, query, instructor.Name, instructor.Email, instructor.Password).Scan(&id)
	return id, err
}

func (r *instructorRepository) GetInstructorByEmail(ctx context.Context, email string) (*models.Instructor, error) {
	query := `SELECT id, name, email, password FROM instructors WHERE email = $1;`

	var instructor models.Instructor
	err := r.db.QueryRow(ctx, query, email).Scan(&instructor.ID, &instructor.Name, &instructor.Email, &instructor.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInstructorNotFound
		}
		return nil, err
	}
	return &instructor, nil
}

func (r *instructorRepository) GetCoursesByInstructor(ctx context.Context, instructorID int64) ([]*models.Course, error) {
	query := `
        SELECT id, name, description, instructor_id
        FROM courses
        WHERE instructor_id = $1;
    `

	rows, err := r.db.Query(ctx, query, instructorID)
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

func (r *instructorRepository) GetInstructorByID(ctx context.Context, id int64) (*models.Instructor, error) {
	query := `SELECT id, name, email, password FROM instructors WHERE id = $1`

	var instructor models.Instructor
	err := r.db.QueryRow(ctx, query, id).Scan(
		&instructor.ID,
		&instructor.Name,
		&instructor.Email,
		&instructor.Password,
	)
	if err != nil {
		return nil, err
	}

	return &instructor, nil
}

func (r *instructorRepository) UpdateInstructor(ctx context.Context, instructor *models.Instructor) error {
	query := `
		UPDATE instructors 
		SET name = $1, email = $2, password = $3 
		WHERE id = $4
	`

	commandTag, err := r.db.Exec(ctx, query, instructor.Name, instructor.Email, instructor.Password, instructor.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("преподаватель не найден")
	}

	return nil
}
