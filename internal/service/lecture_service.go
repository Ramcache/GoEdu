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

func (s *LectureService) GetLecturesByCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.LectureList, error) {
	if req.CourseId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID курса должен быть указан")
	}

	lectures, err := s.lectureRepo.GetLecturesByCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении лекций: %v", err)
	}

	var grpcLectures []*proto.Lecture
	for _, lecture := range lectures {
		grpcLectures = append(grpcLectures, &proto.Lecture{
			Id:       lecture.ID,
			CourseId: lecture.CourseID,
			Title:    lecture.Title,
			Content:  lecture.Content,
		})
	}

	return &proto.LectureList{Lectures: grpcLectures}, nil
}

func (s *LectureService) GetLectureContent(ctx context.Context, req *proto.LectureIDRequest) (*proto.LectureContent, error) {
	if req.LectureId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID лекции должен быть указан")
	}

	lecture, err := s.lectureRepo.GetLectureContent(ctx, req.LectureId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении содержания лекции: %v", err)
	}

	if lecture == nil {
		return nil, status.Errorf(codes.NotFound, "Лекция с ID %d не найдена", req.LectureId)
	}

	return &proto.LectureContent{
		Id:       lecture.ID,
		CourseId: lecture.CourseID,
		Title:    lecture.Title,
		Content:  lecture.Content,
	}, nil
}
