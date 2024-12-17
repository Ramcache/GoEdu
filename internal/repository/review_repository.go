package repository

import (
	"GoEdu/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository interface {
	AddReview(ctx context.Context, studentID, courseID int64, comment string, rating int32) error
	GetReviewsByCourse(ctx context.Context, courseID int64) ([]*models.Review, error)
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

func (r *reviewRepository) GetReviewsByCourse(ctx context.Context, courseID int64) ([]*models.Review, error) {
	query := `
        SELECT id, student_id, course_id, comment, rating, created_at
        FROM reviews
        WHERE course_id = $1;
    `

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.ID, &review.StudentID, &review.CourseID, &review.Comment, &review.Rating, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}
