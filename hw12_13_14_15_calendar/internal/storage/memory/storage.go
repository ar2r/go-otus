package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	items sync.Map
}

func New() *Storage {
	return &Storage{}
}

// Get Вернуть событие по идентификатору.
func (s *Storage) Get(_ context.Context, id uuid.UUID) (*storage.Event, error) {
	if v, exists := s.items.Load(id); exists {
		return v.(*storage.Event), nil
	}
	return nil, storage.ErrNotFound
}

// Add Добавить событие.
func (s *Storage) Add(_ context.Context, event storage.Event) (*storage.Event, error) {
	s.items.Store(event.ID, event)
	v, _ := s.items.Load(event.ID)
	result := v.(storage.Event)
	return &result, nil
}

// Update Обновить событие.
func (s *Storage) Update(_ context.Context, event storage.Event) (*storage.Event, error) {
	if _, exists := s.items.Load(event.ID); !exists {
		return nil, storage.ErrNotFound
	}
	s.items.Store(event.ID, event)
	v, _ := s.items.Load(event.ID)
	result := v.(storage.Event)
	return &result, nil
}

// Delete Удалить событие.
func (s *Storage) Delete(_ context.Context, uuid uuid.UUID) error {
	s.items.Delete(uuid)
	return nil
}

// List Вернуть все события.
func (s *Storage) List(_ context.Context) ([]storage.Event, error) {
	foundItems := make([]storage.Event, 0)
	s.items.Range(func(key, value any) bool {
		foundItems = append(foundItems, value.(storage.Event))
		return true
	})
	return foundItems, nil
}

// FindByDate Найти все события, которые происходят в указанный день.
func (s *Storage) FindByDate(_ context.Context, start time.Time) ([]storage.Event, error) {
	// Обнулить у даты часы минуты секунды
	startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Add(-1 * time.Nanosecond)
	endOfDay := time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, time.Local)

	foundItems := make([]storage.Event, 0)
	s.items.Range(func(key, value any) bool {
		v := value.(storage.Event)
		if v.StartDt.Before(endOfDay) && v.EndDt.After(startOfDay) {
			foundItems = append(foundItems, v)
		}
		return true
	})
	return foundItems, nil
}

// ListByPeriod Найти события, которые пересекается с указанным временным промежутком.
func (s *Storage) ListByPeriod(_ context.Context, startDt time.Time, endDt time.Time) ([]storage.Event, error) {
	startDt = startDt.Add(-1 * time.Nanosecond)
	endDt = endDt.Add(1 * time.Nanosecond)

	foundItems := make([]storage.Event, 0)
	s.items.Range(func(key, value any) bool {
		v := value.(storage.Event)
		if v.StartDt.Before(endDt) && v.EndDt.After(startDt) {
			foundItems = append(foundItems, v)
		}
		return true
	})
	return foundItems, nil
}

// ListByWeek Найти все события, которые происходят в указанной неделе.
// Неделя начинается с понедельника.
func (s *Storage) ListByWeek(ctx context.Context, startDt time.Time) ([]storage.Event, error) {
	startOfWeek := startDt.AddDate(0, 0, -int(startDt.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return s.ListByPeriod(ctx, startOfWeek, endOfWeek)
}

// ListByMonth Найти все события, которые происходят в указанном месяце.
// Месяц начинается с первого числа.
func (s *Storage) ListByMonth(ctx context.Context, startDt time.Time) ([]storage.Event, error) {
	startOfMonth := time.Date(startDt.Year(), startDt.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.ListByPeriod(ctx, startOfMonth, endOfMonth)
}
