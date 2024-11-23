package adapters

import (
	"context"
	"fmt"
	"log/slog"

	memorystorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/memory"
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/pgx"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
)

func New(ctx context.Context, logg *slog.Logger, myConfig *config.Config) (model.EventRepository, error) {
	var eventRepo model.EventRepository
	var err error

	switch myConfig.App.Storage {
	case "memory":
		eventRepo = memorystorage.New()
		logg.Info("Memory adapters initialized")
	case "sql":
		eventRepo, err = sqlstorage.New(ctx, myConfig.Database, logg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SQL storage: %w", err)
		}
		logg.Info("SQL adapters initialized")
	default:
		return nil, fmt.Errorf("invalid adapters type: %s", myConfig.App.Storage)
	}
	return eventRepo, nil
}
