package db

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg config.DatabaseConf, _ *logger.Logger) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s&TimeZone=%s&target_session_attrs=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Database,
		cfg.SSLMode,
		cfg.Timezone,
		cfg.TargetSessionAttr,
	)

	pgxCfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxCfg.ConnConfig.DialFunc = (&net.Dialer{KeepAlive: 15 * time.Second}).DialContext
	pgxCfg.MaxConns = 5
	pgxCfg.MinConns = 3
	pgxCfg.MaxConnLifetime = 5 * time.Minute
	pgxCfg.MaxConnIdleTime = time.Minute

	// todo: Добавить логирование запросов

	db, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create PGX pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PGX pool: %s", err)
	}

	return db, nil
}

func Close(db *pgxpool.Pool) {
	db.Close()
}
