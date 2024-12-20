package service

import (
	"GoEdu/internal/models"
	"GoEdu/internal/repository"
	"GoEdu/proto"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LectureService struct {
	proto.UnimplementedLectureServiceServer
	lectureRepo repository.LectureRepository
	logger      *zap.Logger
}

func NewLectureService(lectureRepo repository.LectureRepository, logger *zap.Logger) *LectureService {
	return &LectureService{
		lectureRepo: lectureRepo,
		logger:      logger,
	}
}

func (s *LectureService) AddLectureToCourse(ctx context.Context, req *proto.LectureRequest) (*proto.Lecture, error) {
	s.logger.Info("Добавление лекции к курсу", zap.Int64("course_id", req.CourseId), zap.String("title", req.Title))

	if req.CourseId == 0 || req.Title == "" || req.Content == "" {
		s.logger.Warn("Некорректные данные для добавления лекции", zap.Int64("course_id", req.CourseId), zap.String("title", req.Title))
		return nil, status.Errorf(codes.InvalidArgument, "Все поля должны быть заполнены")
	}

	lecture := &models.Lecture{
		CourseID: req.CourseId,
		Title:    req.Title,
		Content:  req.Content,
	}

	newLecture, err := s.lectureRepo.AddLectureToCourse(ctx, lecture)
	if err != nil {
		s.logger.Error("Ошибка при добавлении лекции", zap.Error(err), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при добавлении лекции: %v", err)
	}

	s.logger.Info("Лекция успешно добавлена", zap.Int64("lecture_id", newLecture.ID), zap.Int64("course_id", newLecture.CourseID))
	return &proto.Lecture{
		Id:       newLecture.ID,
		CourseId: newLecture.CourseID,
		Title:    newLecture.Title,
		Content:  newLecture.Content,
	}, nil
}

func (s *LectureService) GetLecturesByCourse(ctx context.Context, req *proto.CourseIDRequest) (*proto.LectureList, error) {
	s.logger.Info("Получение лекций для курса", zap.Int64("course_id", req.CourseId))

	if req.CourseId == 0 {
		s.logger.Warn("Некорректный ID курса", zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.InvalidArgument, "ID курса должен быть указан")
	}

	lectures, err := s.lectureRepo.GetLecturesByCourse(ctx, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при получении лекций", zap.Error(err), zap.Int64("course_id", req.CourseId))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении лекций: %v", err)
	}

	var grpcLectures []*proto.Lecture
	for _, lecture := range lectures {
		grpcLectures = append(grpcLectures, &proto.Lecture{
			Id:       lecture.ID,
			CourseId: lecture.CourseID,
			Title:    lecture.Title,
			Content:  lecture.Content,
		})
	}

	s.logger.Info("Лекции успешно получены", zap.Int("count", len(grpcLectures)), zap.Int64("course_id", req.CourseId))
	return &proto.LectureList{Lectures: grpcLectures}, nil
}

func (s *LectureService) GetLectureContent(ctx context.Context, req *proto.LectureIDRequest) (*proto.LectureContent, error) {
	s.logger.Info("Получение содержания лекции", zap.Int64("lecture_id", req.LectureId))

	if req.LectureId == 0 {
		s.logger.Warn("Некорректный ID лекции", zap.Int64("lecture_id", req.LectureId))
		return nil, status.Errorf(codes.InvalidArgument, "ID лекции должен быть указан")
	}

	lecture, err := s.lectureRepo.GetLectureContent(ctx, req.LectureId)
	if err != nil {
		s.logger.Error("Ошибка при получении содержания лекции", zap.Error(err), zap.Int64("lecture_id", req.LectureId))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении содержания лекции: %v", err)
	}

	if lecture == nil {
		s.logger.Warn("Лекция не найдена", zap.Int64("lecture_id", req.LectureId))
		return nil, status.Errorf(codes.NotFound, "Лекция с ID %d не найдена", req.LectureId)
	}

	s.logger.Info("Содержание лекции успешно получено", zap.Int64("lecture_id", lecture.ID))
	return &proto.LectureContent{
		Id:       lecture.ID,
		CourseId: lecture.CourseID,
		Title:    lecture.Title,
		Content:  lecture.Content,
	}, nil
}

func (s *LectureService) UpdateLecture(ctx context.Context, req *proto.UpdateLectureRequest) (*proto.Lecture, error) {
	s.logger.Info("Обновление лекции", zap.Int64("lecture_id", req.Id))

	if req.Id == 0 {
		s.logger.Warn("Некорректный ID лекции", zap.Int64("lecture_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "ID лекции должен быть указан")
	}

	lectureToUpdate := &models.Lecture{
		ID:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	}

	updatedLecture, err := s.lectureRepo.UpdateLecture(ctx, lectureToUpdate)
	if err != nil {
		s.logger.Error("Ошибка при обновлении лекции", zap.Error(err), zap.Int64("lecture_id", req.Id))
		return nil, status.Errorf(codes.Internal, "Ошибка при обновлении лекции: %v", err)
	}

	s.logger.Info("Лекция успешно обновлена", zap.Int64("lecture_id", updatedLecture.ID))
	return &proto.Lecture{
		Id:       updatedLecture.ID,
		CourseId: updatedLecture.CourseID,
		Title:    updatedLecture.Title,
		Content:  updatedLecture.Content,
	}, nil
}

func (s *LectureService) DeleteLecture(ctx context.Context, req *proto.LectureIDRequest) (*proto.Empty, error) {
	s.logger.Info("Удаление лекции", zap.Int64("lecture_id", req.LectureId))

	if req.LectureId == 0 {
		s.logger.Warn("Некорректный ID лекции", zap.Int64("lecture_id", req.LectureId))
		return nil, status.Errorf(codes.InvalidArgument, "ID лекции должен быть указан")
	}

	err := s.lectureRepo.DeleteLecture(ctx, req.LectureId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("Лекция не найдена", zap.Int64("lecture_id", req.LectureId))
			return nil, status.Errorf(codes.NotFound, "Лекция с ID %d не найдена", req.LectureId)
		}
		s.logger.Error("Ошибка при удалении лекции", zap.Error(err), zap.Int64("lecture_id", req.LectureId))
		return nil, status.Errorf(codes.Internal, "Ошибка при удалении лекции: %v", err)
	}

	s.logger.Info("Лекция успешно удалена", zap.Int64("lecture_id", req.LectureId))
	return &proto.Empty{}, nil
}

func (s *LectureService) MarkLectureAsCompleted(ctx context.Context, req *proto.LectureCompletionRequest) (*proto.Empty, error) {
	s.logger.Info("Отметка лекции как завершенной", zap.Int64("lecture_id", req.LectureId), zap.Int64("student_id", req.StudentId))

	if req.StudentId == 0 || req.LectureId == 0 {
		s.logger.Warn("Некорректные данные для отметки завершения лекции", zap.Int64("lecture_id", req.LectureId), zap.Int64("student_id", req.StudentId))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента и ID лекции должны быть указаны")
	}

	err := s.lectureRepo.MarkLectureAsCompleted(ctx, req.StudentId, req.LectureId)
	if err != nil {
		s.logger.Error("Ошибка при отметке лекции как завершенной", zap.Error(err), zap.Int64("lecture_id", req.LectureId), zap.Int64("student_id", req.StudentId))
		return nil, status.Errorf(codes.NotFound, "Ошибка при отметке лекции как завершенной: %v", err)
	}

	s.logger.Info("Лекция успешно отмечена как завершенная", zap.Int64("lecture_id", req.LectureId), zap.Int64("student_id", req.StudentId))
	return &proto.Empty{}, nil
}

func (s *LectureService) GetCourseProgress(ctx context.Context, req *proto.CourseProgressRequest) (*proto.CourseProgress, error) {
	s.logger.Info("Получение прогресса курса", zap.Int64("course_id", req.CourseId), zap.Int64("student_id", req.StudentId))

	if req.StudentId == 0 || req.CourseId == 0 {
		s.logger.Warn("Некорректные данные для получения прогресса курса", zap.Int64("course_id", req.CourseId), zap.Int64("student_id", req.StudentId))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента и ID курса должны быть указаны")
	}

	progress, err := s.lectureRepo.GetCourseProgress(ctx, req.StudentId, req.CourseId)
	if err != nil {
		s.logger.Error("Ошибка при получении прогресса", zap.Error(err), zap.Int64("course_id", req.CourseId), zap.Int64("student_id", req.StudentId))
		return nil, status.Errorf(codes.Internal, "Ошибка при получении прогресса: %v", err)
	}

	s.logger.Info("Прогресс успешно получен", zap.Int64("course_id", req.CourseId), zap.Int64("student_id", req.StudentId), zap.Int("completed_percent", int(progress)))
	return &proto.CourseProgress{
		CourseId:         req.CourseId,
		StudentId:        req.StudentId,
		CompletedPercent: progress,
	}, nil
}

func (s *LectureService) GetRecommendedCourses(ctx context.Context, req *proto.StudentIDRequest) (*proto.CourseList, error) {
	s.logger.Info("Получение рекомендованных курсов", zap.Int64("student_id", req.Id))

	if req.Id == 0 {
		s.logger.Warn("Некорректный ID студента", zap.Int64("student_id", req.Id))
		return nil, status.Errorf(codes.InvalidArgument, "ID студента должен быть указан")
	}

	courses, err := s.lectureRepo.GetRecommendedCourses(ctx, req.Id)
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
