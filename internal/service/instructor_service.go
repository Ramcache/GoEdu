package service

import (
	"GoEdu/internal/config"
	"GoEdu/internal/middleware"
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InstructorService struct {
	proto.UnimplementedInstructorServiceServer
	repo repository.InstructorRepository
	cfg  *config.Config
}

func NewInstructorService(repo repository.InstructorRepository, cfg *config.Config) *InstructorService {
	return &InstructorService{
		repo: repo,
		cfg:  cfg}
}

func (s *InstructorService) RegisterInstructor(ctx context.Context, req *proto.RegisterInstructorRequest) (*proto.Instructor, error) {
	existingInstructor, _ := s.repo.GetInstructorByEmail(ctx, req.Email)
	if existingInstructor != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Преподаватель с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при хэшировании пароля")
	}

	instructor := &models.Instructor{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := s.repo.RegisterInstructor(ctx, instructor)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при регистрации преподавателя: %v", err)
	}

	token, err := middleware.GenerateJWTToken(id, req.Email, "instructor", []byte(s.cfg.JWTSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при создании JWT-токена")
	}

	return &proto.Instructor{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Token: token,
	}, nil
}
