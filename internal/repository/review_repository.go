package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository interface {
	AddReview(ctx context.Context, studentID, courseID int64, comment string, rating int32) error
}

type reviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) AddReview(ctx context.Context, studentID, courseID int64, comment string, rating int32) error {
	query := `
        INSERT INTO reviews (student_id, course_id, comment, rating)
        VALUES ($1, $2, $3, $4);
    `
	_, err := r.db.Exec(ctx, query, studentID, courseID, comment, rating)
	return err
}
