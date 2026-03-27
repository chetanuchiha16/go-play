package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chetanuchiha16/go-play/internal/domain/user"
)

type contextKey string

const user_id contextKey = "user_id"

func (mw MiddlewareManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorisation header required", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorisation format", http.StatusUnauthorized)
			return
		}

		claims, err := user.ValidateToken(parts[1], mw.jwtkey)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return // do not forget to return after hitting errors
		}
		ctx := context.WithValue(r.Context(), user_id, claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
