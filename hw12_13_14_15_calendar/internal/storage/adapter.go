package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	memorystorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage/pgx"
)

func New(ctx context.Context, logg *slog.Logger, myConfig *config.Config) (model.EventRepository, error) {
	var eventRepo model.EventRepository

	switch myConfig.App.Storage {
	case "memory":
		eventRepo = memorystorage.New()
		logg.Info("Memory adapters initialized")
	case "sql":
		dbPool, err := sqlstorage.Connect(ctx, myConfig.Database, logg)
		eventRepo, err = sqlstorage.New(logg, dbPool)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SQL storage: %w", err)
		}
		logg.Info("SQL adapters initialized")
	default:
		return nil, fmt.Errorf("invalid adapters type: %s", myConfig.App.Storage)
	}
	return eventRepo, nil
}
