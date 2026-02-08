package user

import (
	"time"

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

	token := jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	return token.SignedString(jwtkey)
}