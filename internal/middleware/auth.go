package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/deathcore1998/workout-tracker/internal/auth"
)

type contextKey string

const UserID contextKey = "userID"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(writer, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(writer, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := auth.ParseToken(tokenString, secret)
			if err != nil {
				http.Error(writer, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), UserID, claims.UserID)
			next.ServeHTTP(writer, req.WithContext(ctx))
		})
	}
}
