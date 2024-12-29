package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"swimresults-backend/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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

func dbURL() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	return cfg.URL(), nil
}

func Connect(ctx context.Context, migrations fs.FS) (*pgxpool.Pool, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open pool: %w", err)
	}

	url, err := dbURL()
	if err != nil {
		return nil, err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		return nil, fmt.Errorf("migrate new: %s", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}
