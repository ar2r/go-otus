package app

import (
	"context"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model/event"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Logger         Logger
	UserRepository Storage
	PgxPool        *pgxpool.Pool
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, id uuid.UUID, title string) error
	Add(ctx context.Context, event event.Event) (*event.Event, error)
	Update(ctx context.Context, event event.Event) (*event.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]event.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:         logger,
		UserRepository: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id uuid.UUID, title string) error {
	return a.UserRepository.CreateEvent(ctx, id, title)
}
