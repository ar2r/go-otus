package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type Storage struct {
	mu    sync.RWMutex
	items sync.Map
}

func New() *Storage {
	return &Storage{}
}

// CreateEvent Создать событие с проверками на возможные пересечения с другими событиями.
func (s *Storage) CreateEvent(ctx context.Context, userID uuid.UUID, id uuid.UUID, title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := model.Event{
		ID:      id,
		Title:   title,
		StartDt: time.Now(),
		EndDt:   time.Now().Add(time.Hour),
		UserID:  userID,
	}

	foundEvents, err := s.ListByPeriod(ctx, event.StartDt, event.EndDt)
	if err != nil {
		return err
	}

	if len(foundEvents) > 0 {
		return adapters.ErrDateBusy
	}

	_, err = s.Add(ctx, event)
	return err
}

// Get Вернуть событие по идентификатору.
func (s *Storage) Get(_ context.Context, id uuid.UUID) (model.Event, error) {
	if v, exists := s.items.Load(id); exists {
		return v.(model.Event), nil
	}
	return model.Event{}, adapters.ErrNotFound
}

// Add Добавить событие.
func (s *Storage) Add(_ context.Context, e model.Event) (model.Event, error) {
	s.items.Store(e.ID, e)
	v, _ := s.items.Load(e.ID)
	result := v.(model.Event)
	return result, nil
}

// Update Обновить событие.
func (s *Storage) Update(_ context.Context, e model.Event) (model.Event, error) {
	if _, exists := s.items.Load(e.ID); !exists {
		return model.Event{}, adapters.ErrNotFound
	}
	s.items.Store(e.ID, e)
	v, _ := s.items.Load(e.ID)
	result := v.(model.Event)
	return result, nil
}

// Delete Удалить событие.
func (s *Storage) Delete(_ context.Context, id uuid.UUID) error {
	s.items.Delete(id)
	return nil
}

// List Вернуть все события.
func (s *Storage) List(_ context.Context) ([]model.Event, error) {
	foundItems := make([]model.Event, 0)
	s.items.Range(func(_, value any) bool {
		foundItems = append(foundItems, value.(model.Event))
		return true
	})
	return foundItems, nil
}

// ListByDate Найти все события, которые происходят в указанный день.
func (s *Storage) ListByDate(_ context.Context, start time.Time) ([]model.Event, error) {
	// Обнулить у даты часы минуты секунды
	startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Add(-1 * time.Nanosecond)
	endOfDay := time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, time.Local)

	foundItems := make([]model.Event, 0)
	s.items.Range(func(_, value any) bool {
		v := value.(model.Event)
		if v.StartDt.Before(endOfDay) && v.EndDt.After(startOfDay) {
			foundItems = append(foundItems, v)
		}
		return true
	})
	return foundItems, nil
}

// ListByPeriod Найти события, которые пересекается с указанным временным промежутком.
func (s *Storage) ListByPeriod(_ context.Context, startDt time.Time, endDt time.Time) ([]model.Event, error) {
	startDt = startDt.Add(-1 * time.Nanosecond)
	endDt = endDt.Add(1 * time.Nanosecond)

	foundItems := make([]model.Event, 0)
	s.items.Range(func(_, value any) bool {
		v := value.(model.Event)
		if v.StartDt.Before(endDt) && v.EndDt.After(startDt) {
			foundItems = append(foundItems, v)
		}
		return true
	})
	return foundItems, nil
}

// ListByWeek Найти все события, которые происходят в указанной неделе.
// Неделя начинается с понедельника.
func (s *Storage) ListByWeek(ctx context.Context, startDt time.Time) ([]model.Event, error) {
	startOfWeek := startDt.AddDate(0, 0, -int(startDt.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return s.ListByPeriod(ctx, startOfWeek, endOfWeek)
}

// ListByMonth Найти все события, которые происходят в указанном месяце.
// Месяц начинается с первого числа.
func (s *Storage) ListByMonth(ctx context.Context, startDt time.Time) ([]model.Event, error) {
	startOfMonth := time.Date(startDt.Year(), startDt.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.ListByPeriod(ctx, startOfMonth, endOfMonth)
}
