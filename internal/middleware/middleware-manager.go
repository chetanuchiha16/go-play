package middleware

import "github.com/golang-jwt/jwt/v4"

type tokenValidator func(tokenString string) (jwt.MapClaims, error)
type MiddlewareManager struct {
	// jwtkey []byte
	// validateToken tokenValidator
}

func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{
		// jwtkey: jwtkey,
		// validateToken: validateToken,
	}
}
