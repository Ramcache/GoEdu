package service

import (
	"GoEdu/internal/config"
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type StudentService struct {
	proto.UnimplementedStudentServiceServer
	studentRepo repository.StudentRepository
	cfg         config.Config
}

func NewStudentService(studentRepo repository.StudentRepository, cfg config.Config) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		cfg:         cfg,
	}
}

func (s *StudentService) RegisterStudent(ctx context.Context, req *proto.RegisterStudentRequest) (*proto.Student, error) {
	existingStudent, _ := s.studentRepo.GetStudentByEmail(ctx, req.Email)
	if existingStudent != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Студент с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Не удалось хэшировать пароль")
	}

	student := &models.Student{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := s.studentRepo.RegisterStudent(ctx, student)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка регистрации студента: %v", err)
	}

	token, err := s.generateJWTToken(id, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Не удалось создать JWT токен")
	}

	return &proto.Student{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Token: token,
	}, nil
}

func (s *StudentService) generateJWTToken(studentID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    studentID,
		"email": email,
		"exp":   time.Now().Add(time.Duration(s.cfg.TokenExpiryHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecretKey))
}

func (s *StudentService) LoginStudent(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	student, err := s.studentRepo.GetStudentByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Студент с таким email не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Неверный пароль")
	}

	token, err := s.generateJWTToken(student.ID, student.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Не удалось создать JWT-токен")
	}

	return &proto.AuthResponse{
		Id:    student.ID,
		Name:  student.Name,
		Email: student.Email,
		Token: token,
	}, nil
}
