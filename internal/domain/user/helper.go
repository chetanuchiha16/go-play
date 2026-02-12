package user

import (
	"time"
	"fmt"
	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/golang-jwt/jwt/v4"
)
var jwtkey = []byte(config.Load().JWT_SECRET)

func GenerateToken(user_id int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(jwtkey)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v",token.Header["alg"])
		}
		return jwtkey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
} 