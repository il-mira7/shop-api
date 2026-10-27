package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"shopAPI/internal/authserver/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

//go:embed migrations/*.sql
var FS embed.FS

func runMigrations(databaseURL string) error {
	d, err := iofs.New(FS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	return nil
}

func NewPostgresPool(cfg *config.Config, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		return nil, err
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
