package service

import (
	"GoEdu/internal/config"
	"GoEdu/internal/middleware"
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InstructorService struct {
	proto.UnimplementedInstructorServiceServer
	repo   repository.InstructorRepository
	cfg    *config.Config
	logger *zap.Logger
}

func NewInstructorService(repo repository.InstructorRepository, cfg *config.Config, logger *zap.Logger) *InstructorService {
	return &InstructorService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *InstructorService) RegisterInstructor(ctx context.Context, req *proto.RegisterInstructorRequest) (*proto.Instructor, error) {
	existingInstructor, err := s.repo.GetInstructorByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Ошибка при проверке существующего email", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при проверке существующего email: %v", err)
	}

	if existingInstructor != nil {
		s.logger.Warn("Преподаватель с таким email уже существует", zap.String("email", req.Email))
		return nil, status.Errorf(codes.AlreadyExists, "Преподаватель с таким email уже существует")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Ошибка при хэшировании пароля", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при хэшировании пароля")
	}

	instructor := &models.Instructor{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := s.repo.RegisterInstructor(ctx, instructor)
	if err != nil {
		s.logger.Error("Ошибка при регистрации преподавателя", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при регистрации преподавателя: %v", err)
	}

	token, err := middleware.GenerateJWTToken(id, req.Email, "instructor", []byte(s.cfg.JWTSecretKey), s.cfg.TokenExpiryHours, s.logger)
	if err != nil {
		s.logger.Error("Не удалось создать JWT-токен", zap.Error(err), zap.Int64("instructor_id", id))
		return nil, status.Errorf(codes.Internal, "Не удалось создать JWT-токен")
	}

	s.logger.Info("Преподаватель успешно зарегистрирован", zap.Int64("instructor_id", id), zap.String("email", req.Email))

	return &proto.Instructor{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Token: token,
	}, nil
}

func (s *InstructorService) GetCoursesByInstructor(ctx context.Context, req *proto.InstructorIDRequest) (*proto.CourseList, error) {
	s.logger.Info("Получение курсов для преподавателя", zap.Int64("instructor_id", req.InstructorId))

	if req.InstructorId == 0 {
		s.logger.Warn("Некорректный ID преподавателя", zap.Int64("instructor_id", req.InstructorId))
		return nil, status.Errorf(codes.InvalidArgument, "ID преподавателя должен быть указан")
	}

	courses, err := s.repo.GetCoursesByInstructor(ctx, req.InstructorId)
	if err != nil {
		s.logger.Error("Ошибка при получении курсов", zap.Error(err), zap.Int64("instructor_id", req.InstructorId))
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

	s.logger.Info("Курсы успешно получены", zap.Int("count", len(grpcCourses)), zap.Int64("instructor_id", req.InstructorId))
	return &proto.CourseList{Courses: grpcCourses}, nil
}
