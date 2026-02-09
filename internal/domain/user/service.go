package user

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/chetanuchiha16/go-play/db"
)

type Service interface {
	CreateUser(ctx context.Context, args db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, id int64) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]db.User, error)
	Login(ctx context.Context, email, password string) (user db.User, token string, err error)
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
	password_hash, err := bcrypt.GenerateFromPassword([]byte(args.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("error making hash")
		return db.User{}, err
	}
	// args = db.CreateUserParams{
	// 	Name : args.Name,
	// 	Email: args.Email,
	// 	PasswordHash : string(password_hash),

	// }
	args.PasswordHash = string(password_hash)
	return s.store.CreateUser(ctx, args)
}

func (s *userService) GetUser(ctx context.Context, id int64) (db.User, error) {
	return s.store.GetUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.store.DeleteUser(ctx, id)
}

// handler expects something that implements the servide interface
// making userService implement it
func (s *userService) ListUsers(ctx context.Context) ([]db.User, error) {
	return s.store.ListUsers(ctx)
}

func (s *userService) Login(ctx context.Context, email, password string) (user db.User, token string, err error) {
	user, err = s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return db.User{}, "", err
	}
	token, err = GenerateToken(user.ID)
	if err != nil {
		return db.User{}, "", err
	}
	return user, token, nil
}
