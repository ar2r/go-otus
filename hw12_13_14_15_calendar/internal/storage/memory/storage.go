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

type StartEndDt struct {
	StartDt time.Time //Дата и время события;
	EndDt   time.Time //Длительность события (или дата и время окончания);
}

type EventId struct {
	Id int //Идентификатор события;
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

func (s *Storage) Update(value interface{}) error {
	// todo: Переписать на замену. Это временное рабочее решение.
	if eventIdObject, ok := value.(EventId); ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, exists := s.items[eventIdObject.Id]; !exists {
			return ErrNotFound
		}
		s.items[eventIdObject.Id] = value
		return nil
	}
	return ErrNotFound
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

// FindByWeek Найти все события, которые происходят в указанной неделе.
// Неделя начинается с понедельника.
func (s *Storage) FindByWeek(startDt time.Time) map[int]interface{} {
	startOfWeek := startDt.AddDate(0, 0, -int(startDt.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return s.FindByPeriod(startOfWeek, endOfWeek)
}

// FindByMonth Найти все события, которые происходят в указанном месяце.
// Месяц начинается с первого числа.
func (s *Storage) FindByMonth(startDt time.Time) map[int]interface{} {
	startOfMonth := time.Date(startDt.Year(), startDt.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.FindByPeriod(startOfMonth, endOfMonth)
}
