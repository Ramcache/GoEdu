package service

import (
	"GoEdu/internal/config"
	"GoEdu/internal/middleware"
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
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

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nameRegex  = regexp.MustCompile(`^[\p{L}\p{N}\s'-]+$`)
)

func (s *InstructorService) RegisterInstructor(ctx context.Context, req *proto.RegisterInstructorRequest) (*proto.Instructor, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Имя преподавателя не может быть пустым")
	}
	if len(req.Email) > 255 {
		return nil, status.Errorf(codes.InvalidArgument, "Email не должен превышать 255 символов")
	}
	if len(req.Name) < 2 || len(req.Name) > 255 {
		return nil, status.Errorf(codes.InvalidArgument, "Имя должно содержать от 2 до 255 символов")
	}
	if !nameRegex.MatchString(req.Name) {
		return nil, status.Errorf(codes.InvalidArgument, "Имя может содержать только буквы, цифры, пробелы, дефисы и апострофы")
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Email преподавателя не может быть пустым")
	}
	if !emailRegex.MatchString(req.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "Некорректный формат email")
	}
	if len(req.Password) < 6 {
		return nil, status.Errorf(codes.InvalidArgument, "Пароль должен содержать не менее 6 символов")
	}

	existingInstructor, err := s.repo.GetInstructorByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repository.ErrInstructorNotFound) {
		s.logger.Error("Ошибка при проверке существующего email", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при проверке существующего email")
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
		return nil, status.Errorf(codes.Internal, "Ошибка при регистрации преподавателя")
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

func (s *InstructorService) LoginInstructor(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	instructor, err := s.repo.GetInstructorByEmail(ctx, req.Email)
	if err != nil {
		// Если преподаватель не найден, возвращаем codes.NotFound
		if errors.Is(err, repository.ErrInstructorNotFound) {
			return nil, status.Errorf(codes.NotFound, "Преподаватель не найден")
		}
		// Для любых других ошибок возвращаем codes.Internal
		return nil, status.Errorf(codes.Internal, "Ошибка при получении данных преподавателя: %v", err)
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(instructor.Password), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Неверный пароль")
	}

	// Генерируем JWT-токен
	token, err := middleware.GenerateJWTToken(instructor.ID, instructor.Email, "instructor", []byte(s.cfg.JWTSecretKey), 24, s.logger)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка генерации токена")
	}

	// Логируем успешный вход
	s.logger.Info("Профиль преподавателя успешно получен", zap.String("Email", req.Email))
	return &proto.AuthResponse{
		Id:    instructor.ID,
		Name:  instructor.Name,
		Email: instructor.Email,
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

func (s *InstructorService) UpdateInstructor(ctx context.Context, req *proto.UpdateInstructorRequest) (*proto.Instructor, error) {
	instructor, err := s.repo.GetInstructorByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Преподаватель не найден", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Преподаватель не найден")
	}

	// Обновление профиля, если данные переданы
	if req.Name != "" {
		instructor.Name = req.Name
	}
	if req.Email != "" {
		instructor.Email = req.Email
	}

	// Обновление пароля, если передан новый пароль
	if req.NewPassword != "" {
		// Проверка текущего пароля
		if err := bcrypt.CompareHashAndPassword([]byte(instructor.Password), []byte(req.CurrentPassword)); err != nil {
			s.logger.Warn("Неверный текущий пароль")
			return nil, status.Errorf(codes.InvalidArgument, "Неверный текущий пароль")
		}

		// Хэширование нового пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error("Ошибка при хэшировании пароля", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Ошибка при изменении пароля")
		}
		instructor.Password = string(hashedPassword)
	}

	// Сохранение изменений в базе данных
	if err := s.repo.UpdateInstructor(ctx, instructor); err != nil {
		s.logger.Error("Ошибка при обновлении данных преподавателя", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при обновлении данных")
	}

	return instructor.ToProto(), nil
}

func (s *InstructorService) GetInstructorByID(ctx context.Context, req *proto.GetInstructorRequest) (*proto.Instructor, error) {
	instructor, err := s.repo.GetInstructorByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Преподаватель не найден", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Преподаватель не найден")
	}

	return instructor.ToProto(), nil
}
