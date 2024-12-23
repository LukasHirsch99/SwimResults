package database

import (
	"context"
	"fmt"
	"swimresults-backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func loadConfig() (*pgxpool.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	))

}

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := loadConfig()
	if err != nil {
    return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
