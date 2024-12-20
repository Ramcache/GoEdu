package repository

import (
	"GoEdu/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *models.Course, tx pgx.Tx) (int, error)
	GetAllCourses(ctx context.Context) ([]*models.Course, error)
	GetCourseByID(ctx context.Context, id int64) (*models.Course, error)
	UpdateCourse(ctx context.Context, tx pgx.Tx, id int64, name, description string) (*models.Course, error)
	DeleteCourse(ctx context.Context, id int64) (bool, error)
	GetCoursesByInstructor(ctx context.Context, instructorID int64) ([]*models.Course, error)
	SearchCourses(ctx context.Context, keyword string) ([]*models.Course, error)
	GetRecommendedCourses(ctx context.Context, studentID int64) ([]*models.Course, error)
}

type courseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) CreateCourse(ctx context.Context, course *models.Course, tx pgx.Tx) (int, error) {
	query := `
        INSERT INTO courses (name, description, instructor_id)
        VALUES ($1, $2, $3)
        RETURNING id;
    `
	var id int
	err := tx.QueryRow(ctx, query, course.Name, course.Description, course.InstructorID).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return 0, fmt.Errorf("преподаватель с ID %d не существует", course.InstructorID)
			case "23505":
				return 0, fmt.Errorf("курс с таким названием уже существует")
			}
		}
		return 0, err
	}

	return id, nil
}

func (r *courseRepository) GetAllCourses(ctx context.Context) ([]*models.Course, error) {
	query := `SELECT id, name, description, instructor_id FROM courses;`

	rows, err := r.db.Query(ctx, query)
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

	log.Printf("Количество курсов в базе: %d", len(courses))

	return courses, nil
}

func (r *courseRepository) GetCourseByID(ctx context.Context, id int64) (*models.Course, error) {
	query := `SELECT id, name, description FROM courses WHERE id = $1;`

	var course models.Course
	err := r.db.QueryRow(ctx, query, id).Scan(&course.ID, &course.Name, &course.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) UpdateCourse(ctx context.Context, tx pgx.Tx, id int64, name, description string) (*models.Course, error) {
	var existingID int64
	queryCheck := "SELECT id FROM courses WHERE name = $1 AND id != $2"
	err := tx.QueryRow(ctx, queryCheck, name, id).Scan(&existingID)

	if err == nil {
		return nil, fmt.Errorf("duplicate name: %w", pgx.ErrNoRows)
	} else if err != pgx.ErrNoRows {
		return nil, err
	}

	queryUpdate := `
        UPDATE courses
        SET name = $1, description = $2
        WHERE id = $3
        RETURNING id, name, description;
    `

	var course models.Course
	err = tx.QueryRow(ctx, queryUpdate, name, description, id).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
	)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) DeleteCourse(ctx context.Context, id int64) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid course ID: %d", id)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM courses WHERE id = $1;`
	commandTag, err := tx.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	if commandTag.RowsAffected() == 0 {
		return false, nil
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *courseRepository) GetCoursesByInstructor(ctx context.Context, instructorID int64) ([]*models.Course, error) {
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

func (r *courseRepository) SearchCourses(ctx context.Context, keyword string) ([]*models.Course, error) {
	query := `
        SELECT id, name, description, instructor_id
        FROM courses
        WHERE name ILIKE $1 OR description ILIKE $1;
    `

	rows, err := r.db.Query(ctx, query, "%"+keyword+"%")
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

func (r *courseRepository) GetRecommendedCourses(ctx context.Context, studentID int64) ([]*models.Course, error) {
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
