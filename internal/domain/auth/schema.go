package auth

import "github.com/chetanuchiha16/go-play/internal/domain/user"


type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
type LoginResponse struct {
	Token string       `json:"token"`
	User  user.UserResponse `json:"user"`
}
