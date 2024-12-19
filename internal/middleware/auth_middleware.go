package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(jwtSecretKey []byte, logger *zap.Logger) grpc.UnaryServerInterceptor {
	whitelist := map[string]bool{
		"/education.StudentService/RegisterStudent":       true,
		"/education.StudentService/LoginStudent":          true,
		"/education.InstructorService/RegisterInstructor": true,

		"/education.EducationService/GetCourses":      true,
		"/education.EducationService/GetCourseByID":   true,
		"/education.EducationService/SearchCourses":   true,
		"/education.ReviewService/GetReviewsByCourse": true,

		"/education.LectureService/GetLecturesByCourse": true,
		"/education.LectureService/GetLectureContent":   true,
	}

	roleProtectedMethods := map[string]string{
		"/education.EducationService/CreateCourseByInstructor": "instructor",
		"/education.StudentService/GetStudentProfile":          "student",
		"/education.StudentService/UpdateStudentProfile":       "student",
		"/education.EnrollmentService/EnrollStudent":           "student",
		"/education.EnrollmentService/GetStudentsByCourse":     "instructor",
		"/education.EnrollmentService/GetCoursesByStudent":     "student",
		"/education.EnrollmentService/UnEnrollStudent":         "student",
		"/education.LectureService/UpdateLecture":              "instructor",
		"/education.LectureService/DeleteLecture":              "instructor",
		"/education.LectureService/MarkLectureAsCompleted":     "student",
		"/education.LectureService/GetCourseProgress":          "student",
		"/education.LectureService/GetRecommendedCourses":      "student",
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Info("Запрос метода", zap.String("method", info.FullMethod))

		if whitelist[info.FullMethod] {
			logger.Info("Метод находится в whitelist, пропуск аутентификации", zap.String("method", info.FullMethod))
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Warn("Метаданные отсутствуют", zap.String("method", info.FullMethod))
			return nil, status.Errorf(codes.Unauthenticated, "Метаданные отсутствуют")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			logger.Warn("Токен отсутствует", zap.String("method", info.FullMethod))
			return nil, status.Errorf(codes.Unauthenticated, "Токен отсутствует")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("Некорректный токен", zap.String("method", info.FullMethod))
				return nil, status.Errorf(codes.Unauthenticated, "Некорректный токен")
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			logger.Warn("Недействительный токен", zap.Error(err), zap.String("method", info.FullMethod))
			return nil, status.Errorf(codes.Unauthenticated, "Недействительный токен")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("Невозможно получить claims из токена", zap.String("method", info.FullMethod))
			return nil, status.Errorf(codes.Unauthenticated, "Невозможно получить claims из токена")
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			logger.Warn("Роль пользователя отсутствует в токене", zap.String("method", info.FullMethod))
			return nil, status.Errorf(codes.PermissionDenied, "Роль пользователя отсутствует в токене")
		}

		if requiredRole, exists := roleProtectedMethods[info.FullMethod]; exists {
			if userRole != requiredRole {
				logger.Warn("Роль пользователя не соответствует требованиям метода", zap.String("method", info.FullMethod), zap.String("required_role", requiredRole), zap.String("user_role", userRole))
				return nil, status.Errorf(codes.PermissionDenied, "Доступ запрещён для вашей роли: требуется %s", requiredRole)
			}
		}

		logger.Info("Аутентификация успешна", zap.String("method", info.FullMethod), zap.String("user_role", userRole))
		return handler(ctx, req)
	}
}

func GenerateJWTToken(userID int64, email, role string, secretKey []byte, tokenExpirationHours int, logger *zap.Logger) (string, error) {
	expirationTime := time.Now().Add(time.Duration(tokenExpirationHours) * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		logger.Error("Ошибка генерации JWT-токена", zap.Error(err), zap.Int64("user_id", userID), zap.String("email", email))
		return "", err
	}

	logger.Info("JWT-токен успешно сгенерирован", zap.Int64("user_id", userID), zap.String("email", email), zap.String("role", role), zap.Time("expires_at", expirationTime))
	return signedToken, nil
}
