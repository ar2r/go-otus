package app

import (
	"context"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, id, title string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.Storage.CreateEvent(ctx, id, title)
}
