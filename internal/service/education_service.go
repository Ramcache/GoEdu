package service

import (
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type EducationService struct {
	proto.UnimplementedEducationServiceServer
	courseRepo repository.CourseRepository
}

func NewEducationService(courseRepo repository.CourseRepository) *EducationService {
	return &EducationService{
		courseRepo: courseRepo,
	}
}

func (s *EducationService) CreateCourse(ctx context.Context, req *proto.NewCourseRequest) (*proto.Course, error) {
	course := &models.Course{
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := s.courseRepo.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}

	return &proto.Course{
		Id:          strconv.FormatInt(int64(id), 10),
		Name:        course.Name,
		Description: course.Description,
	}, nil
}

func (s *EducationService) GetCourses(ctx context.Context, req *proto.Empty) (*proto.CourseList, error) {
	courses, err := s.courseRepo.GetAllCourses(ctx)
	if err != nil {
		return nil, err
	}

	var grpcCourses []*proto.Course
	for _, c := range courses {
		grpcCourses = append(grpcCourses, &proto.Course{
			Id:          strconv.Itoa(c.ID),
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return &proto.CourseList{Courses: grpcCourses}, nil
}

func (s *EducationService) GetCourseByID(ctx context.Context, req *proto.CourseIDRequest) (*proto.Course, error) {
	course, err := s.courseRepo.GetCourseByID(ctx, req.CourseId)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, status.Errorf(codes.NotFound, "Курс с ID %d не найден", req.CourseId)
	}

	return &proto.Course{
		Id:          strconv.FormatInt(int64(course.ID), 10),
		Name:        course.Name,
		Description: course.Description,
	}, nil
}

func (s *EducationService) UpdateCourse(ctx context.Context, req *proto.UpdateCourseRequest) (*proto.Course, error) {
	updatedCourse, err := s.courseRepo.UpdateCourse(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	if updatedCourse == nil {
		return nil, status.Errorf(codes.NotFound, "Курс с ID %d не найден", req.Id)
	}

	return &proto.Course{
		Id:          strconv.FormatInt(int64(updatedCourse.ID), 10),
		Name:        updatedCourse.Name,
		Description: updatedCourse.Description,
	}, nil
}

func (s *EducationService) DeleteCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.Empty, error) {
	deleted, err := s.courseRepo.DeleteCourse(ctx, req.CourseId)
	if err != nil {
		return nil, err
	}

	if !deleted {
		return nil, status.Errorf(codes.NotFound, "Курс с ID %d не найден", req.CourseId)
	}

	return &proto.Empty{}, nil
}
