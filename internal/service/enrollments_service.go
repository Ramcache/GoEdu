package service

import (
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EnrollmentService struct {
	proto.UnimplementedEnrollmentServiceServer
	enrollmentRepo repository.EnrollmentRepository
	logger         *zap.Logger
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, logger *zap.Logger) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
		logger:         logger,
	}
}

func (s *EnrollmentService) EnrollStudent(ctx context.Context, req *proto.EnrollmentRequest) (*proto.Empty, error) {
	s.logger.Info("Запись студента на курс", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))

	if req.StudentId == 0 || req.CourseId == 0 {
		s.logger.Warn("Некорректные данные для записи на курс", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента и курса должны быть указаны")
	}

	err := s.enrollmentRepo.EnrollStudent(ctx, req.StudentId, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при записи студента на курс", zap.Error(err), zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при записи студента на курс: %v", err)
	}

	s.logger.Info("Студент успешно записан на курс", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
	return &proto.Empty{}, nil
}

func (s *EnrollmentService) GetStudentsByCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.StudentList, error) {
	s.logger.Info("Получение студентов по курсу", zap.Int64("course_id", req.CourseId))

	if req.CourseId == 0 {
		s.logger.Warn("Некорректный ID курса", zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.InvalidArgument, "ID курса должен быть указан")
	}

	students, err := s.enrollmentRepo.GetStudentsByCourse(ctx, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при получении студентов", zap.Error(err), zap.Int64("course_id", req.CourseId))
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

	s.logger.Info("Студенты успешно получены", zap.Int("count", len(grpcStudents)), zap.Int64("course_id", req.CourseId))
	return &proto.StudentList{Students: grpcStudents}, nil
}

func (s *EnrollmentService) UnEnrollStudent(ctx context.Context, req *proto.UnEnrollRequest) (*proto.Empty, error) {
	s.logger.Info("Удаление студента с курса", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))

	if req.StudentId == 0 || req.CourseId == 0 {
		s.logger.Warn("Некорректные данные для удаления с курса", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента и курса должны быть указаны")
	}

	err := s.enrollmentRepo.UnEnrollStudent(ctx, req.StudentId, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при удалении студента с курса", zap.Error(err), zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при удалении студента с курса: %v", err)
	}

	s.logger.Info("Студент успешно удалён с курса", zap.Int64("student_id", req.StudentId), zap.Int64("course_id", req.CourseId))
	return &proto.Empty{}, nil
}

func (s *EnrollmentService) GetCoursesByStudent(ctx context.Context, req *proto.StudentIDRequest) (*proto.CourseList, error) {
	s.logger.Info("Получение курсов для студента", zap.Int64("student_id", req.Id))

	if req.Id == 0 {
		s.logger.Warn("Некорректный ID студента", zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента должен быть указан")
	}

	courses, err := s.enrollmentRepo.GetCoursesByStudent(ctx, req.Id)
	if err != nil {
		s.logger.Error("Ошибка при получении курсов", zap.Error(err), zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении курсов: %v", err)
	}

	var grpcCourses []*proto.Course
	for _, course := range courses {
		grpcCourses = append(grpcCourses, &proto.Course{
			Id:          course.ID,
			Name:        course.Name,
			Description: course.Description,
		})
	}

	s.logger.Info("Курсы успешно получены", zap.Int("count", len(grpcCourses)), zap.Int64("student_id", req.Id))
	return &proto.CourseList{Courses: grpcCourses}, nil
}
