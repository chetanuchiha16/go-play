package auth

import "github.com/chetanuchiha16/go-play/internal/domain/user"


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string       `json:"token"`
	User  user.UserResponse `json:"user"`
}
