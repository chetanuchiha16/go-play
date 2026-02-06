package user

import (
	"context"

	"github.com/chetanuchiha16/go-play/db"
)

type Service interface {
	CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, id int64) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]db.User, error)
}

type userService struct {
	store Store
}

func NewService(s Store) *userService {
	return &userService{
		store: s,
	}
}

func (s *userService) CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error) {
	
	return s.store.CreateUser(ctx, args)
}

func (s *userService) GetUser(ctx context.Context, id int64) (db.User, error) {
	return s.store.GetUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.store.DeleteUser(ctx, id)
}

//handler expects something that implements the servide interface
// making userService implement it
func (s *userService) ListUsers(ctx context.Context) ([]db.User, error) {
	return s.store.ListUsers(ctx)
}
