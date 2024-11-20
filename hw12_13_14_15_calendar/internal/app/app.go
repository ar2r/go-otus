package app

import (
	"context"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type App struct {
	Logger         *slog.Logger
	userRepository model.EventRepository
}

func New(logg *slog.Logger, repo model.EventRepository) *App {
	return &App{
		Logger:         logg,
		userRepository: repo,
	}
}

func (a *App) CreateEvent(ctx context.Context, e model.Event) error {
	if _, err := a.userRepository.Add(ctx, e); err != nil {
		return err
	}

	return nil
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (model.Event, error) {
	return a.userRepository.Get(ctx, id)
}

func (a *App) UpdateEvent(ctx context.Context, ev model.Event) (model.Event, error) {
	return a.userRepository.Update(ctx, ev)
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.userRepository.Delete(ctx, id)
}
