package memorystorage

import (
	"errors"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		name      string
		operation func(s *Storage) error
		expected  map[int]interface{}
	}{
		{
			name: "Add item",
			operation: func(s *Storage) error {
				s.Add("item1")
				return nil
			},
			expected: map[int]interface{}{1: "item1"},
		},
		{
			name: "Add multiple items",
			operation: func(s *Storage) error {
				s.Add("item1")
				s.Add("item2")
				return nil
			},
			expected: map[int]interface{}{
				1: "item1",
				2: "item2",
			},
		},
		{
			name: "Delete item",
			operation: func(s *Storage) error {
				id, _ := s.Add("item1")
				return s.Delete(id)
			},
			expected: map[int]interface{}{},
		},
		{
			name: "Delete non-existent item",
			operation: func(s *Storage) error {
				return s.Delete(999)
			},
			expected: map[int]interface{}{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := New()
			err := tt.operation(storage)
			if err != nil && !errors.Is(err, ErrNotFound) {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := storage.List(); !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestAddIntersectingEvent(t *testing.T) {
	storage := New()
	// 2000-01-01 12:00:00 +0000 UTC
	start1 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end1 := start1.Add(2 * time.Hour)
	_, err := storage.Add(StartEndDt{StartDt: start1, EndDt: end1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2000-01-01 13:00:00 +0000 UTC
	start2 := start1.Add(1 * time.Hour)
	// 2000-01-01 15:00:00 +0000 UTC
	end2 := start2.Add(2 * time.Hour)
	_, err = storage.Add(StartEndDt{StartDt: start2, EndDt: end2})
	if err != ErrIntersect {
		t.Errorf("expected error %v, got %v", ErrIntersect, err)
	}
}

func TestFindByDate(t *testing.T) {
	storage := New()
	// 2000-01-01 12:00:00 +0000 UTC
	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)
	storage.Add(StartEndDt{StartDt: start, EndDt: end})

	tests := []struct {
		name     string
		date     time.Time
		expected map[int]interface{}
	}{
		{
			name: "Find existing event",
			date: start,
			expected: map[int]interface{}{
				1: StartEndDt{StartDt: start, EndDt: end},
			},
		},
		{
			name:     "Find non-existing event",
			date:     start.Add(24 * time.Hour),
			expected: map[int]interface{}{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := storage.FindByDate(tt.date)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestFindByPeriod(t *testing.T) {
	storage := New()
	// 2000-01-01 12:00:00 +0000 UTC
	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)
	storage.Add(StartEndDt{StartDt: start, EndDt: end})

	tests := []struct {
		name     string
		startDt  time.Time
		endDt    time.Time
		expected map[int]interface{}
	}{
		{
			name:    "Find overlapping event",
			startDt: start.Add(-1 * time.Hour),
			endDt:   start.Add(1 * time.Hour),
			expected: map[int]interface{}{
				1: StartEndDt{StartDt: start, EndDt: end},
			},
		},
		{
			name:     "Find non-overlapping event",
			startDt:  start.Add(3 * time.Hour),
			endDt:    start.Add(4 * time.Hour),
			expected: map[int]interface{}{},
		},
		{
			name:    "Find event within range",
			startDt: start,
			endDt:   end,
			expected: map[int]interface{}{
				1: StartEndDt{StartDt: start, EndDt: end},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := storage.FindByPeriod(tt.startDt, tt.endDt)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func equal(a, b map[int]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
