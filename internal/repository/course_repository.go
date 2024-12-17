package repository

import (
	"GoEdu/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *models.Course) (int, error)
	GetAllCourses(ctx context.Context) ([]*models.Course, error)
	GetCourseByID(ctx context.Context, id int64) (*models.Course, error)
	UpdateCourse(ctx context.Context, id int64, name, description string) (*models.Course, error)
	DeleteCourse(ctx context.Context, id int64) (bool, error)
	GetCoursesByInstructor(ctx context.Context, instructorID int64) ([]*models.Course, error)
	SearchCourses(ctx context.Context, keyword string) ([]*models.Course, error)
}

type courseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) CreateCourse(ctx context.Context, course *models.Course) (int, error) {
	query := `
        INSERT INTO courses (name, description)
        VALUES ($1, $2)
        RETURNING id;
    `
	var id int
	err := r.db.QueryRow(ctx, query, course.Name, course.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *courseRepository) GetAllCourses(ctx context.Context) ([]*models.Course, error) {
	query := `SELECT id, name, description FROM courses;`
	rows, err := r.db.Query(ctx, query)
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

func (r *courseRepository) UpdateCourse(ctx context.Context, id int64, name, description string) (*models.Course, error) {
	query := `
        UPDATE courses
        SET name = $1, description = $2
        WHERE id = $3
        RETURNING id, name, description;
    `

	var course models.Course
	err := r.db.QueryRow(ctx, query, name, description, id).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) DeleteCourse(ctx context.Context, id int64) (bool, error) {
	query := `DELETE FROM courses WHERE id = $1;`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	if commandTag.RowsAffected() == 0 {
		return false, nil
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
