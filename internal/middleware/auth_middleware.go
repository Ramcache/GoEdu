package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(jwtSecretKey []byte) grpc.UnaryServerInterceptor {
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
		if whitelist[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Метаданные отсутствуют")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Токен отсутствует")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, status.Errorf(codes.Unauthenticated, "Некорректный токен")
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "Недействительный токен")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Невозможно получить claims из токена")
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "Роль пользователя отсутствует в токене")
		}

		if requiredRole, exists := roleProtectedMethods[info.FullMethod]; exists {
			if userRole != requiredRole {
				return nil, status.Errorf(codes.PermissionDenied, "Доступ запрещён для вашей роли: требуется %s", requiredRole)
			}
		}

		return handler(ctx, req)
	}
}

func GenerateJWTToken(userID int64, email, role string, secretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
