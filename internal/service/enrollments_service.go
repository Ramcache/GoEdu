package service

import (
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EnrollmentService struct {
	proto.UnimplementedEnrollmentServiceServer
	enrollmentRepo repository.EnrollmentRepository
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
	}
}

func (s *EnrollmentService) EnrollStudent(ctx context.Context, req *proto.EnrollmentRequest) (*proto.Empty, error) {
	if req.StudentId == 0 || req.CourseId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID студента и курса должны быть указаны")
	}

	err := s.enrollmentRepo.EnrollStudent(ctx, req.StudentId, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при записи студента на курс: %v", err)
	}

	return &proto.Empty{}, nil
}

func (s *EnrollmentService) GetStudentsByCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.StudentList, error) {
	if req.CourseId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID курса должен быть указан")
	}

	students, err := s.enrollmentRepo.GetStudentsByCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении студентов: %v", err)
	}

	var grpcStudents []*proto.Student
	for _, student := range students {
		grpcStudents = append(grpcStudents, &proto.Student{
			Id:    student.ID,
			Name:  student.Name,
			Email: student.Email,
			Token: "",
		})
	}

	return &proto.StudentList{Students: grpcStudents}, nil
}

func (s *EnrollmentService) UnEnrollStudent(ctx context.Context, req *proto.EnrollmentRequest) (*proto.Empty, error) {
	if req.StudentId == 0 || req.CourseId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ID студента и курса должны быть указаны")
	}

	err := s.enrollmentRepo.UnEnrollStudent(ctx, req.StudentId, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при удалении студента с курса: %v", err)
	}

	return &proto.Empty{}, nil
}
