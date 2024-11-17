package main

import (
	"errors"
	"fmt"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/golang-migrate/migrate/v4"
	// Register the pgx database driver.
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	// Register the file source for migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func MigrateRun(logg Logger, conf config.DatabaseConf, rerun bool) error {
	logg.Info("Connecting to database...")
	dsn := fmt.Sprintf(
		"pgx://%s:%s@%s/%s?sslmode=%s&TimeZone=%s&target_session_attrs=%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Database,
		conf.SSLMode,
		conf.Timezone,
		conf.TargetSessionAttr,
	)

	logg.Info("Loading migrations...")
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("migration: error while create connection: %w", err)
	}

	if rerun {
		logg.Info("Rollback migrations...")
		if err = m.Down(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migration: %w", err)
			}
			logg.Info(fmt.Sprintf("migration: %s", err))
		}
	}

	logg.Info("Run migrations...")
	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration: %w", err)
		}
		logg.Info(fmt.Sprintf("migration: %s", err))
	}

	logg.Info("Closing database connection...")
	s, e := m.Close()
	if s != nil {
		return fmt.Errorf("error while closing migration: %w", e)
	}
	if e != nil {
		return fmt.Errorf("error while closing migration: %w", e)
	}
	return nil
}
