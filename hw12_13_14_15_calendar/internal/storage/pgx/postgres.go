package sqlstorage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, conf database.Config, logg *slog.Logger) (*pgxpool.Pool, error) {
	pgxPool, err := database.Connect(ctx, conf)
	if err != nil {
		logg.Error(fmt.Sprintf("failed to create connection to database: %s", err))
		return nil, err
	}
	return pgxPool, nil
}
