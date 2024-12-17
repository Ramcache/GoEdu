package service

import (
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		Id:          int64(id),
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
			Id:          c.ID,
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
		Id:          course.ID,
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
		Id:          updatedCourse.ID,
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

func (s *EducationService) CreateCourseByInstructor(ctx context.Context, req *proto.InstructorCourseRequest) (*proto.Course, error) {
	if req.InstructorId == 0 || req.Name == "" || req.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID преподавателя, название и описание курса должны быть указаны")
	}

	course := &models.Course{
		Name:         req.Name,
		Description:  req.Description,
		InstructorID: req.InstructorId,
	}

	courseID, err := s.courseRepo.CreateCourse(ctx, course)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при создании курса: %v", err)
	}

	return &proto.Course{
		Id:           int64(courseID),
		Name:         req.Name,
		Description:  req.Description,
		InstructorId: req.InstructorId,
	}, nil
}

func (s *EducationService) GetCoursesByInstructor(ctx context.Context, req *proto.InstructorIDRequest) (*proto.CourseList, error) {
	if req.InstructorId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID преподавателя должен быть указан")
	}

	courses, err := s.courseRepo.GetCoursesByInstructor(ctx, req.InstructorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении курсов: %v", err)
	}

	var grpcCourses []*proto.Course
	for _, course := range courses {
		grpcCourses = append(grpcCourses, &proto.Course{
			Id:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			InstructorId: course.InstructorID,
		})
	}

	return &proto.CourseList{Courses: grpcCourses}, nil
}
