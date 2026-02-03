package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	// pool, err := pgxpool.New(ctx, connStr)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, config)

}
