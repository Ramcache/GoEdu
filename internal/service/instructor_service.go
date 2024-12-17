package service

import (
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
}

func NewInstructorService(repo repository.InstructorRepository) *InstructorService {
	return &InstructorService{repo: repo}
}

// RegisterInstructor регистрирует нового преподавателя
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

	return &proto.Instructor{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}, nil
}
