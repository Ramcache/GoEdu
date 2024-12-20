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

type StudentService struct {
	proto.UnimplementedStudentServiceServer
	studentRepo repository.StudentRepository
	cfg         *config.Config
	logger      *zap.Logger
}

func NewStudentService(studentRepo repository.StudentRepository, cfg *config.Config, logger *zap.Logger) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		cfg:         cfg,
		logger:      logger,
	}
}

func (s *StudentService) RegisterStudent(ctx context.Context, req *proto.RegisterStudentRequest) (*proto.Student, error) {
	s.logger.Info("Регистрация студента началась", zap.String("email", req.Email))

	existingStudent, _ := s.studentRepo.GetStudentByEmail(ctx, req.Email)
	if existingStudent != nil {
		s.logger.Warn("Студент с таким email уже существует", zap.String("email", req.Email))
		return nil, status.Errorf(codes.AlreadyExists, "Студент с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Ошибка хэширования пароля", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Не удалось хэшировать пароль")
	}

	student := &models.Student{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := s.studentRepo.RegisterStudent(ctx, student)
	if err != nil {
		s.logger.Error("Ошибка регистрации студента", zap.Error(err), zap.String("email", req.Email))
		return nil, status.Errorf(codes.Internal, "Ошибка регистрации студента: %v", err)
	}

	token, err := middleware.GenerateJWTToken(id, req.Email, "student", []byte(s.cfg.JWTSecretKey), s.cfg.TokenExpiryHours, s.logger)
	if err != nil {
		s.logger.Error("Ошибка генерации JWT токена", zap.Error(err), zap.Int64("student_id", id))
		return nil, status.Errorf(codes.Internal, "Не удалось создать JWT токен")
	}

	s.logger.Info("Студент успешно зарегистрирован", zap.Int64("student_id", id), zap.String("email", req.Email))
	return &proto.Student{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Token: token,
	}, nil
}

func (s *StudentService) LoginStudent(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	s.logger.Info("Авторизация студента началась", zap.String("email", req.Email))

	student, err := s.studentRepo.GetStudentByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("Студент не найден", zap.String("email", req.Email))
		return nil, status.Errorf(codes.NotFound, "Студент с таким email не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password))
	if err != nil {
		// Обработка ошибки, если пароль неверный
		return nil, status.Errorf(codes.Unauthenticated, "Неверный пароль")
	}

	token, err := middleware.GenerateJWTToken(student.ID, student.Email, "student", []byte(s.cfg.JWTSecretKey), s.cfg.TokenExpiryHours, s.logger)
	if err != nil {
		s.logger.Error("Ошибка генерации JWT токена", zap.Error(err), zap.Int64("student_id", student.ID))
		return nil, status.Errorf(codes.Internal, "Не удалось создать JWT токен")
	}

	s.logger.Info("Студент успешно авторизован", zap.Int64("student_id", student.ID), zap.String("email", req.Email))
	return &proto.AuthResponse{
		Id:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Token: token,
	}, nil
}

func (s *StudentService) GetStudentProfile(ctx context.Context, req *proto.StudentIDRequest) (*proto.Student, error) {
	s.logger.Info("Получение профиля студента", zap.Int64("student_id", req.Id))

	student, err := s.studentRepo.GetStudentByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Ошибка получения профиля студента", zap.Error(err), zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении профиля студента: %v", err)
	}

	if student == nil {
		s.logger.Warn("Студент не найден", zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.NotFound, "Студент с ID %d не найден", req.Id)
	}

	s.logger.Info("Профиль студента успешно получен", zap.Int64("student_id", req.Id))
	return &proto.Student{
		Id:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Token: "",
	}, nil
}

func (s *StudentService) UpdateStudentProfile(ctx context.Context, req *proto.UpdateStudentRequest) (*proto.Student, error) {
	s.logger.Info("Обновление профиля студента", zap.Int64("student_id", req.Id))

	existingStudent, err := s.studentRepo.GetStudentByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Ошибка получения студента", zap.Error(err), zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении студента: %v", err)
	}

	if existingStudent == nil {
		s.logger.Warn("Студент не найден", zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.NotFound, "Студент с ID %d не найден", req.Id)
	}

	if req.Name != "" {
		existingStudent.Name = req.Name
	}
	if req.Email != "" {
		existingStudent.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error("Ошибка хэширования пароля", zap.Error(err), zap.Int64("student_id", req.Id))
			return nil, status.Errorf(codes.Internal, "Не удалось хэшировать пароль")
		}
		existingStudent.Password = string(hashedPassword)
	}

	updatedStudent, err := s.studentRepo.UpdateStudent(ctx, existingStudent)
	if err != nil {
		s.logger.Error("Ошибка обновления профиля", zap.Error(err), zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при обновлении профиля: %v", err)
	}

	s.logger.Info("Профиль студента успешно обновлён", zap.Int64("student_id", updatedStudent.ID))
	return &proto.Student{
		Id:    updatedStudent.ID,
		Name:  updatedStudent.Name,
		Email: updatedStudent.Email,
		Token: "",
	}, nil
}
