package service

import (
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LectureService struct {
	proto.UnimplementedLectureServiceServer
	lectureRepo repository.LectureRepository
}

func NewLectureService(lectureRepo repository.LectureRepository) *LectureService {
	return &LectureService{lectureRepo: lectureRepo}
}

func (s *LectureService) AddLectureToCourse(ctx context.Context, req *proto.LectureRequest) (*proto.Lecture, error) {
	if req.CourseId == 0 || req.Title == "" || req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Все поля должны быть заполнены")
	}

	lecture := &models.Lecture{
		CourseID: req.CourseId,
		Title:    req.Title,
		Content:  req.Content,
	}

	newLecture, err := s.lectureRepo.AddLectureToCourse(ctx, lecture)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при добавлении лекции: %v", err)
	}

	return &proto.Lecture{
		Id:       newLecture.ID,
		CourseId: newLecture.CourseID,
		Title:    newLecture.Title,
		Content:  newLecture.Content,
	}, nil
}
