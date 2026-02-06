package user

import (
	"context"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	GetUser(ctx context.Context, id int64) (db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}

type SQLStore struct {
	pool        *pgxpool.Pool
	*db.Queries // get create del user are methods of *db.queries they are now belongs to Mystore ?
}

func NewStore(pool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		pool:    pool,
		Queries: db.New(pool),
	}
}
