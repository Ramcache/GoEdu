package service

import (
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type EducationService struct {
	proto.UnimplementedEducationServiceServer
	db         *pgxpool.Pool
	courseRepo repository.CourseRepository
	logger     *zap.Logger
}

func NewEducationService(db *pgxpool.Pool, courseRepo repository.CourseRepository, logger *zap.Logger) *EducationService {
	return &EducationService{
		db:         db,
		courseRepo: courseRepo,
		logger:     logger,
	}
}

func (s *EducationService) CreateCourse(ctx context.Context, req *proto.NewCourseRequest) (*proto.Course, error) {
	s.logger.Info("Создание курса", zap.String("name", req.Name), zap.Int64("instructor_id", req.InstructorId))

	if req.InstructorId == 0 {
		s.logger.Warn("Некорректный ID преподавателя")
		return nil, status.Errorf(codes.InvalidArgument, "ID преподавателя должен быть указан")
	}

	if req.Name == "" {
		s.logger.Warn("Пустое название курса")
		return nil, status.Errorf(codes.InvalidArgument, "Название курса не может быть пустым")
	}

	if req.Description == "" {
		s.logger.Warn("Пустое описание курса")
		return nil, status.Errorf(codes.InvalidArgument, "Описание курса не может быть пустым")
	}

	course := &models.Course{
		Name:         req.Name,
		Description:  req.Description,
		InstructorID: req.InstructorId,
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		s.logger.Error("Не удалось начать транзакцию", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при начале транзакции: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	id, err := s.courseRepo.CreateCourse(ctx, course, tx)
	if err != nil {
		if strings.Contains(err.Error(), "курс с таким названием уже существует") {
			s.logger.Warn("Курс с таким названием уже существует", zap.String("name", req.Name))
			return nil, status.Errorf(codes.AlreadyExists, "Курс с таким названием уже существует")
		}
		if strings.Contains(err.Error(), fmt.Sprintf("преподаватель с ID %d не существует", course.InstructorID)) {
			s.logger.Warn("Некорректный ID преподавателя", zap.Int64("instructor_id", req.InstructorId))
			return nil, status.Errorf(codes.InvalidArgument, "Преподаватель с ID %d не существует", req.InstructorId)
		}

		s.logger.Error("Ошибка при создании курса", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при создании курса")
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("Не удалось зафиксировать транзакцию", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при фиксации транзакции: %v", err)
	}

	s.logger.Info("Курс успешно создан", zap.Int64("course_id", int64(id)))
	return &proto.Course{
		Id:           int64(id),
		Name:         course.Name,
		Description:  course.Description,
		InstructorId: course.InstructorID,
	}, nil
}

func (s *EducationService) GetCourses(ctx context.Context, req *proto.Empty) (*proto.CourseList, error) {
	s.logger.Info("Получение всех курсов")

	courses, err := s.courseRepo.GetAllCourses(ctx)
	if err != nil {
		s.logger.Error("Ошибка при получении курсов", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении курсов: %v", err)
	}

	s.logger.Info("Курсы успешно получены", zap.Int("count", len(courses)))

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
	s.logger.Info("Получение курса по ID", zap.Int64("course_id", req.CourseId))

	course, err := s.courseRepo.GetCourseByID(ctx, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при получении курса", zap.Error(err), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении курса: %v", err)
	}

	if course == nil {
		s.logger.Warn("Курс не найден", zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.NotFound, "Курс с ID %d не найден", req.CourseId)
	}

	s.logger.Info("Курс успешно получен", zap.Int64("course_id", req.CourseId))
	return &proto.Course{
		Id:          course.ID,
		Name:        course.Name,
		Description: course.Description,
	}, nil
}

func (s *EducationService) UpdateCourse(ctx context.Context, req *proto.UpdateCourseRequest) (*proto.Course, error) {
	if req.Id <= 0 {
		s.logger.Warn("Некорректный ID курса", zap.Int64("course_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "Некорректный ID курса")
	}
	if len(req.Name) > 255 {
		s.logger.Warn("Слишком длинное имя курса", zap.Int64("course_id", req.Id), zap.String("name", req.Name))
		return nil, status.Errorf(codes.InvalidArgument, "Слишком длинное имя курса")
	}
	if req.Name == "" {
		s.logger.Warn("Пустое имя курса", zap.Int64("course_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "Имя курса не может быть пустым")
	}
	if req.Description == "" {
		s.logger.Warn("Пустое описание курса", zap.Int64("course_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "Описание курса не может быть пустым")
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		s.logger.Error("Не удалось начать транзакцию", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Не удалось начать транзакцию: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			s.logger.Error("Транзакция откатилась из-за ошибки", zap.Error(err))
		} else {
			tx.Commit(ctx)
			s.logger.Info("Транзакция успешно завершена")
		}
	}()

	course, err := s.courseRepo.GetCourseByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Ошибка при получении курса", zap.Error(err), zap.Int64("course_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении курса: %v", err)
	}
	if course == nil {
		s.logger.Warn("Курс не найден", zap.Int64("course_id", req.Id))
		return nil, status.Errorf(codes.NotFound, "Курс с ID %d не найден", req.Id)
	}

	updatedCourse, err := s.courseRepo.UpdateCourse(ctx, tx, req.Id, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("Курс с таким названием уже существует", zap.String("name", req.Name))
			return nil, status.Errorf(codes.AlreadyExists, "Курс с таким названием уже существует")
		}
		s.logger.Error("Ошибка при обновлении курса", zap.Error(err), zap.Int64("course_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при обновлении курса: %v", err)
	}

	s.logger.Info("Курс успешно обновлен", zap.Int64("course_id", updatedCourse.ID))
	return &proto.Course{
		Id:          updatedCourse.ID,
		Name:        updatedCourse.Name,
		Description: updatedCourse.Description,
	}, nil
}

func (s *EducationService) DeleteCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.Empty, error) {
	s.logger.Info("Удаление курса", zap.Int64("course_id", req.CourseId))

	deleted, err := s.courseRepo.DeleteCourse(ctx, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при удалении курса", zap.Error(err), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при удалении курса: %v", err)
	}

	if !deleted {
		s.logger.Warn("Курс не найден", zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.NotFound, "Курс с ID %d не найден", req.CourseId)
	}

	s.logger.Info("Курс успешно удален", zap.Int64("course_id", req.CourseId))
	return &proto.Empty{}, nil
}

func (s *EducationService) GetCoursesByInstructor(ctx context.Context, req *proto.InstructorIDRequest) (*proto.CourseList, error) {
	s.logger.Info("Получение курсов для преподавателя", zap.Int64("instructor_id", req.InstructorId))

	if req.InstructorId == 0 {
		s.logger.Warn("Некорректный ID преподавателя", zap.Int64("instructor_id", req.InstructorId))
		return nil, status.Errorf(codes.InvalidArgument, "ID преподавателя должен быть указан")
	}

	courses, err := s.courseRepo.GetCoursesByInstructor(ctx, req.InstructorId)
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

func (s *EducationService) SearchCourses(ctx context.Context, req *proto.SearchRequest) (*proto.CourseList, error) {
	s.logger.Info("Поиск курсов", zap.String("keyword", req.Keyword))

	if req.Keyword == "" {
		s.logger.Warn("Пустое ключевое слово для поиска")
		return nil, status.Errorf(codes.InvalidArgument, "Ключевое слово для поиска не может быть пустым")
	}

	courses, err := s.courseRepo.SearchCourses(ctx, req.Keyword)
	if err != nil {
		s.logger.Error("Ошибка при поиске курсов", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Ошибка при поиске курсов: %v", err)
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

	s.logger.Info("Курсы успешно найдены", zap.Int("count", len(grpcCourses)))
	return &proto.CourseList{Courses: grpcCourses}, nil
}

func (s *EducationService) GetRecommendedCourses(ctx context.Context, req *proto.StudentIDRequest) (*proto.CourseList, error) {
	s.logger.Info("Получение рекомендованных курсов", zap.Int64("student_id", req.Id))

	if req.Id == 0 {
		s.logger.Warn("Некорректный ID студента", zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента должен быть указан")
	}

	courses, err := s.courseRepo.GetRecommendedCourses(ctx, req.Id)
	if err != nil {
		s.logger.Error("Ошибка при получении рекомендованных курсов", zap.Error(err), zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении рекомендованных курсов: %v", err)
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

	s.logger.Info("Рекомендованные курсы успешно получены", zap.Int("count", len(grpcCourses)), zap.Int64("student_id", req.Id))
	return &proto.CourseList{Courses: grpcCourses}, nil
}
