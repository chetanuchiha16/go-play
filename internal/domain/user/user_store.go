package user

import (
	"context"

	"github.com/chetanuchiha16/go-play/db"
)

type UserStore interface {
	GetUser(ctx context.Context, id int64) (db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit int32) ([]db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}
