package database

import (
	// "context"

	"context"
	"fmt"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	// GetUser(ctx context.Context, id int64) (db.User, error)
	// CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	// DeleteUser(ctx context.Context, id int64) error
	// ListUsers(ctx context.Context) ([]db.User, error)
	// GetUserByEmail(ctx context.Context, email string) (db.User, error)

	db.Querier
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

func (s *SQLStore) ExecTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	// Create a new Query instance using the transaction
	q := db.New(tx) 
	err = fn(q)
	
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}