package service

import (
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	proto.UnimplementedReviewServiceServer
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) AddReviewToCourse(ctx context.Context, req *proto.ReviewRequest) (*proto.Empty, error) {
	if req.StudentId == 0 || req.CourseId == 0 || req.Rating < 1 || req.Rating > 5 {
		return nil, status.Errorf(codes.InvalidArgument, "Некорректные данные: ID студента, ID курса и оценка должны быть указаны корректно")
	}

	err := s.reviewRepo.AddReview(ctx, req.StudentId, req.CourseId, req.Comment, req.Rating)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при добавлении отзыва: %v", err)
	}

	return &proto.Empty{}, nil
}

func (s *ReviewService) GetReviewsByCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.ReviewList, error) {
	if req.CourseId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID курса должен быть указан")
	}

	reviews, err := s.reviewRepo.GetReviewsByCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении отзывов: %v", err)
	}

	var grpcReviews []*proto.Review
	for _, review := range reviews {
		grpcReviews = append(grpcReviews, &proto.Review{
			Id:        review.ID,
			StudentId: review.StudentID,
			CourseId:  review.CourseID,
			Comment:   review.Comment,
			Rating:    review.Rating,
			CreatedAt: review.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &proto.ReviewList{Reviews: grpcReviews}, nil
}