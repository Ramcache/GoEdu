package repository

import (
	"GoEdu/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository interface {
	RegisterStudent(ctx context.Context, student *models.Student) (int64, error)
	GetStudentByEmail(ctx context.Context, email string) (*models.Student, error)
	GetStudentByID(ctx context.Context, id int64) (*models.Student, error)
}

type studentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) RegisterStudent(ctx context.Context, student *models.Student) (int64, error) {
	query := `
        INSERT INTO students (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id;
    `
	var id int64
	err := r.db.QueryRow(ctx, query, student.Name, student.Email, student.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *studentRepository) GetStudentByEmail(ctx context.Context, email string) (*models.Student, error) {
	query := `SELECT id, name, email, password FROM students WHERE email = $1;`

	var student models.Student
	err := r.db.QueryRow(ctx, query, email).Scan(&student.ID, &student.Name, &student.Email, &student.Password)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *studentRepository) GetStudentByID(ctx context.Context, id int64) (*models.Student, error) {
	query := `SELECT id, name, email FROM students WHERE id = $1;`

	var student models.Student
	err := r.db.QueryRow(ctx, query, id).Scan(&student.ID, &student.Name, &student.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &student, nil
}
