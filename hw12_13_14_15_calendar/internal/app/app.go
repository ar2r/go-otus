package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	Logger  Logger
	Storage Storage
	PgxPool *pgxpool.Pool
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Add(ctx context.Context, event storage.Event) (*storage.Event, error)
	Update(ctx context.Context, event storage.Event) (*storage.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	//return a.Storage.CreateEvent(ctx, id, title)
	return storage.ErrNotImplemented
}
