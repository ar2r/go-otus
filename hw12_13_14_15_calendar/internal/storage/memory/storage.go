package memorystorage

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("item not found")
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

func (s *Storage) Add(value interface{}) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.items[id] = value
	s.nextID++
	return id
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

func (s *Storage) FindByDate(start time.Time) map[int]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	foundItems := make(map[int]interface{})
	for k, v := range s.items {
		if se, ok := v.(StartEndDt); ok {
			if se.StartDt.Before(start) && se.EndDt.After(start) {
				foundItems[k] = v
			}
		}
	}
	return foundItems
}
