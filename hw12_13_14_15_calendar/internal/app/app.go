package app

import (
	"context"

	dto2 "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/dto"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
)

//go:generate mockgen -source=app.go -destination=mocks/app.go -package=mocks

type App struct {
	repository model.EventRepository
}

type IApplication interface {
	CreateEvent(ctx context.Context, dto dto2.CreateEventDto) (model.Event, error)
	UpdateEvent(ctx context.Context, dto dto2.UpdateEventDto) (model.Event, error)
	DeleteEvent(ctx context.Context, dto dto2.DeleteEventDto) error
	ListByDate(ctx context.Context, dto dto2.ListByDateDto) ([]model.Event, error)
	ListByWeek(ctx context.Context, dto dto2.ListByDateDto) ([]model.Event, error)
	ListByMonth(ctx context.Context, dto dto2.ListByDateDto) ([]model.Event, error)
}

func New(repo model.EventRepository) *App {
	return &App{
		repository: repo,
	}
}

// CreateEvent Создание события.
func (a *App) CreateEvent(ctx context.Context, dto dto2.CreateEventDto) (model.Event, error) {
	m := dto.ToModel()
	return a.repository.Add(ctx, m)
}

// UpdateEvent Обновление события.
func (a *App) UpdateEvent(ctx context.Context, dto dto2.UpdateEventDto) (model.Event, error) {
	m := dto.ToModel()
	return a.repository.Update(ctx, m)
}

// DeleteEvent Удаление события.
func (a *App) DeleteEvent(ctx context.Context, dto dto2.DeleteEventDto) error {
	return a.repository.Delete(ctx, dto.ID)
}

// ListByDate Получение списка событий на дату.
func (a *App) ListByDate(ctx context.Context, dto dto2.ListByDateDto) ([]model.Event, error) {
	return a.repository.ListByDate(ctx, dto.Date)
}

// ListByWeek Получение списка событий на неделю.
func (a *App) ListByWeek(ctx context.Context, dto dto2.ListByDateDto) ([]model.Event, error) {
	return a.repository.ListByWeek(ctx, dto.Date)
}

// ListByMonth Получение списка событий за месяц.
func (a *App) ListByMonth(ctx context.Context, dto dto2.ListByDateDto) ([]model.Event, error) {
	return a.repository.ListByMonth(ctx, dto.Date)
}
