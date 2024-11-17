package app

import (
	"context"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model/event"
	"github.com/google/uuid"
)

type App struct {
	Logger         Logger
	userRepository event.Repository
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	InfoRaw(msg string)
	Warn(msg string)
	Error(msg string)
}

func New(logger Logger, repo event.Repository) *App {
	return &App{
		Logger:         logger,
		userRepository: repo,
	}
}

func (a *App) CreateEvent(ctx context.Context, userID, id uuid.UUID, title string) error {
	return a.userRepository.CreateEvent(ctx, userID, id, title)
}
