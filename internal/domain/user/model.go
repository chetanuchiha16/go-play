package user

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Name string `json:"name"`
}

type CreateUserShema struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
