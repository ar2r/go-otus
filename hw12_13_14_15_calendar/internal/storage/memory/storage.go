package memorystorage

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("event not found")
	ErrDateBusy = errors.New("intersecting events")
)

type Storage struct {
	mu     sync.RWMutex
	items  map[int]interface{}
	nextID int
}

func New() *Storage {
	return &Storage{
		items:  make(map[int]interface{}),
		nextID: 1,
	}
}

func (s *Storage) Add(value interface{}) (int, error) {
	// Проверить, что нет событий, которые пересекаются с новым событием
	if startEndObject, ok := value.(StartEndDt); ok {
		foundItems := s.FindByPeriod(startEndObject.StartDt, startEndObject.EndDt)
		if len(foundItems) > 0 {
			return -1, ErrDateBusy
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.items[id] = value
	s.nextID++
	return id, nil
}

func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[id]; !exists {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}

func (s *Storage) List() map[int]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	foundItems := make(map[int]interface{})
	for k, v := range s.items {
		foundItems[k] = v
	}
	return foundItems
}

type StartEndDt struct {
	StartDt time.Time //Дата и время события;
	EndDt   time.Time //Длительность события (или дата и время окончания);
}

// FindByDate Найти все события, которые происходят в указанный день.
func (s *Storage) FindByDate(start time.Time) map[int]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Обнулить у даты часы минуты секунды
	startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Add(-1 * time.Nanosecond)
	endOfDay := time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, time.Local)

	foundItems := make(map[int]interface{})
	for k, v := range s.items {
		if startEndObject, ok := v.(StartEndDt); ok {
			if startEndObject.StartDt.Before(endOfDay) && startEndObject.EndDt.After(startOfDay) {
				foundItems[k] = v
			}
		}
	}
	return foundItems
}

// FindByPeriod Найти события, которые пересекается с указанным временным промежутком.
func (s *Storage) FindByPeriod(startDt time.Time, endDt time.Time) map[int]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	startDt = startDt.Add(-1 * time.Nanosecond)
	endDt = endDt.Add(1 * time.Nanosecond)

	foundItems := make(map[int]interface{})
	for k, v := range s.items {
		if startEndObject, ok := v.(StartEndDt); ok {
			if startEndObject.StartDt.Before(endDt) && startEndObject.EndDt.After(startDt) {
				foundItems[k] = v
			}
		}
	}
	return foundItems
}
