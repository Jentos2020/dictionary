package store

import (
	"context"
	"leetgo/config"
	"leetgo/internal/errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	migrationsPath := "file://" + cfg.App.MigrationsPath
	m, err := migrate.New(
		migrationsPath,
		cfg.GetMigrateDsn(),
	)
	if err != nil {
		return errors.NewF("failed to init migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.NewF("migration up error: %w", err)
	}

	log.Info("Migrations applied")
	return nil
}
