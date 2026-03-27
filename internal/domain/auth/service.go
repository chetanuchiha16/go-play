package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (user db.User, token string, err error)
	GenerateToken(user_id int64) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type authService struct {
	store  user.UserStore
	jwtkey []byte
}

func NewAuthService(store user.UserStore, jwtkey []byte) *authService {
	return &authService{
		store:  store,
		jwtkey: jwtkey,
	}
}
func (s *authService) Login(ctx context.Context, email, password string) (user db.User, token string, err error) {
	user, err = s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return db.User{}, "", err
	}
	token, err = s.GenerateToken(user.ID)
	if err != nil {
		return db.User{}, "", err
	}
	return user, token, nil
}

func (s authService) GenerateToken(user_id int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtkey)
}

func (s authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %w, %v", errors.ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return s.jwtkey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidToken
}
