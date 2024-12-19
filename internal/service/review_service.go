package service

import (
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	proto.UnimplementedReviewServiceServer
	reviewRepo repository.ReviewRepository
	logger     *zap.Logger
}

func NewReviewService(reviewRepo repository.ReviewRepository, logger *zap.Logger) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		logger:     logger,
	}
}

func (s *ReviewService) AddReviewToCourse(ctx context.Context, req *proto.ReviewRequest) (*proto.Empty, error) {
	s.logger.Info("Добавление отзыва к курсу", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId), zap.Int("rating", int(req.Rating)))

	if req.StudentId == 0 || req.CourseId == 0 || req.Rating < 1 || req.Rating > 5 {
		s.logger.Warn("Некорректные данные для отзыва", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId), zap.Int("rating", int(req.Rating)))
		return nil, status.Errorf(codes.InvalidArgument, "Некорректные данные: ID студента, ID курса и оценка должны быть указаны корректно")
	}

	err := s.reviewRepo.AddReview(ctx, req.StudentId, req.CourseId, req.Comment, req.Rating)
	if err != nil {
		s.logger.Error("Ошибка при добавлении отзыва", zap.Error(err), zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при добавлении отзыва: %v", err)
	}

	s.logger.Info("Отзыв успешно добавлен", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
	return &proto.Empty{}, nil
}

func (s *ReviewService) GetReviewsByCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.ReviewList, error) {
	s.logger.Info("Получение отзывов для курса", zap.Int64("course_id", req.CourseId))

	if req.CourseId == 0 {
		s.logger.Warn("Некорректный ID курса", zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.InvalidArgument, "ID курса должен быть указан")
	}

	reviews, err := s.reviewRepo.GetReviewsByCourse(ctx, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при получении отзывов", zap.Error(err), zap.Int64("course_id", req.CourseId))
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

	s.logger.Info("Отзывы успешно получены", zap.Int("count", len(grpcReviews)), zap.Int64("course_id", req.CourseId))
	return &proto.ReviewList{Reviews: grpcReviews}, nil
}
