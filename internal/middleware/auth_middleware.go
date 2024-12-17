package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var jwtSecretKey = []byte("your_jwt_secret_key")

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
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

		return handler(ctx, req)
	}
}
