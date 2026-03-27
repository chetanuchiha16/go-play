package user

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/chetanuchiha16/go-play/db"
)

type UserService interface {
	CreateUser(ctx context.Context, args CreateUserShema) (db.User, error)
	GetUser(ctx context.Context, id int64) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit int32) ([]db.User, error)
	// Login(ctx context.Context, email, password string) (user db.User, token string, err error)
}

type userService struct {
	store UserStore
	// jwtkey []byte
}

func NewUserService(s UserStore) *userService {
	return &userService{
		store: s,
		// jwtkey: jwtkey,
	}
}

func (s *userService) CreateUser(ctx context.Context, args CreateUserShema) (db.User, error) {
	password_hash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error making hash")
		return db.User{}, err
	}
	

	user := db.CreateUserParams{Name: args.Name, PasswordHash: string(password_hash), Email: args.Email}
	return s.store.CreateUser(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id int64) (db.User, error) {
	return s.store.GetUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.store.DeleteUser(ctx, id)
}

// handler expects something that implements the servide interface
// making userService implement it
func (s *userService) ListUsers(ctx context.Context, limit int32) ([]db.User, error) {
	return s.store.ListUsers(ctx, limit)
}


