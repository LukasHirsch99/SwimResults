package database

import (
	"context"
	"swimresults-backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), cfg.URL())
	// db, err := pgx.Connect(context.Background(), cfg.URL())
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
