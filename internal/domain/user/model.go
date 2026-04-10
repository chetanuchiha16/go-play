package user

import (
	"github.com/chetanuchiha16/go-play/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Name string `json:"name"`
}

type CreateUserShema struct {
	Name     string `json:"name" validate:"required,min=3,max=100" form:"name"`
	Email    string `json:"email" validate:"required,email" form:"email"`
	Password string `json:"password" validate:"required,min=8,max=72" form:"password"`
}

type UserResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func NewUserResponse(user db.User) UserResponse {
	return UserResponse{
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
	}
}
