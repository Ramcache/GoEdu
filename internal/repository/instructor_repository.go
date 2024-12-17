package repository

import (
	"GoEdu/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InstructorRepository interface {
	RegisterInstructor(ctx context.Context, instructor *models.Instructor) (int64, error)
	GetInstructorByEmail(ctx context.Context, email string) (*models.Instructor, error)
}

type instructorRepository struct {
	db *pgxpool.Pool
}

func NewInstructorRepository(db *pgxpool.Pool) InstructorRepository {
	return &instructorRepository{db: db}
}

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
		return nil, err
	}
	return &instructor, nil
}
