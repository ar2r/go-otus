package memorystorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

var ctx = context.Background()

func equal(a, b []storage.Event) bool {
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

func TestStorage(t *testing.T) {

	event := createStubEvent("item1", time.Time{}, time.Time{})

	tests := []struct {
		name      string
		operation func(s *Storage) error
		expected  []storage.Event
	}{
		{
			name: "Add item",
			operation: func(s *Storage) error {
				s.Add(ctx, event)
				return nil
			},
			expected: []storage.Event{event},
		},
		{
			name: "Delete item",
			operation: func(s *Storage) error {
				s.Add(ctx, event)
				return s.Delete(ctx, event.ID)
			},
			expected: []storage.Event{},
		},
		{
			name: "Delete non-existent item",
			operation: func(s *Storage) error {
				id, _ := uuid.NewV7()
				return s.Delete(ctx, id)
			},
			expected: []storage.Event{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			memStorage := New()
			err := tt.operation(memStorage)
			if err != nil && !errors.Is(err, storage.ErrNotFound) {
				t.Fatalf("unexpected error: %v", err)
			}
			if got, _ := memStorage.List(ctx); !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	memStorage := New()
	event := createStubEvent("item1", time.Time{}, time.Time{})
	memStorage.Add(ctx, event)

	event.Title = "item2"
	_, err := memStorage.Update(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateNotExistentEvent(t *testing.T) {
	memStorage := New()
	event := createStubEvent("item1", time.Time{}, time.Time{})
	_, err := memStorage.Update(ctx, event)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected error %v, got %v", storage.ErrNotFound, err)
	}
}

func TestUpdate(t *testing.T) {
	memStorage := New()
	// 2000-01-01 12:00:00 +0000 UTC
	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)
	event1 := createStubEvent("item1", start, end)

	event2 := event1
	event2.Title = "new title"

	tests := []struct {
		name      string
		operation func() error
		expected  []storage.Event
	}{
		{
			name: "Update item",
			operation: func() error {
				memStorage = New()
				memStorage.Add(ctx, event1)
				_, err := memStorage.Update(ctx, event2)
				return err
			},
			expected: []storage.Event{event2},
		},
		{
			name: "Update non-existent item",
			operation: func() error {
				memStorage = New()
				event := createStubEvent("not exist event", start, end)
				_, err := memStorage.Update(ctx, event)
				return err
			},
			expected: []storage.Event{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.operation()
			if err != nil && !errors.Is(err, storage.ErrNotFound) {
				t.Fatalf("unexpected error: %v", err)
			}
			if got, _ := memStorage.List(ctx); !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

//func TestFindByDate(t *testing.T) {
//	storage := New()
//	// 2000-01-01 12:00:00 +0000 UTC
//	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
//	// 2000-01-01 14:00:00 +0000 UTC
//	end := start.Add(2 * time.Hour)
//	storage.Add(StartEndDt{StartDt: start, EndDt: end})
//
//	tests := []struct {
//		name     string
//		date     time.Time
//		expected map[int]interface{}
//	}{
//		{
//			name: "Find existing event",
//			date: start,
//			expected: map[int]interface{}{
//				1: StartEndDt{StartDt: start, EndDt: end},
//			},
//		},
//		{
//			name:     "Find non-existing event",
//			date:     start.Add(24 * time.Hour),
//			expected: map[int]interface{}{},
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//			got := storage.FindByDate(tt.date)
//			if !equal(got, tt.expected) {
//				t.Errorf("expected %v, got %v", tt.expected, got)
//			}
//		})
//	}
//}
//
//func TestFindByPeriod(t *testing.T) {
//	storage := New()
//	// 2000-01-01 12:00:00 +0000 UTC
//	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
//	// 2000-01-01 14:00:00 +0000 UTC
//	end := start.Add(2 * time.Hour)
//	storage.Add(StartEndDt{StartDt: start, EndDt: end})
//
//	tests := []struct {
//		name     string
//		startDt  time.Time
//		endDt    time.Time
//		expected map[int]interface{}
//	}{
//		{
//			name:    "Find overlapping event",
//			startDt: start.Add(-1 * time.Hour),
//			endDt:   start.Add(1 * time.Hour),
//			expected: map[int]interface{}{
//				1: StartEndDt{StartDt: start, EndDt: end},
//			},
//		},
//		{
//			name:     "Find non-overlapping event",
//			startDt:  start.Add(3 * time.Hour),
//			endDt:    start.Add(4 * time.Hour),
//			expected: map[int]interface{}{},
//		},
//		{
//			name:    "Find event within range",
//			startDt: start,
//			endDt:   end,
//			expected: map[int]interface{}{
//				1: StartEndDt{StartDt: start, EndDt: end},
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//			got := storage.ListByPeriod(tt.startDt, tt.endDt)
//			if !equal(got, tt.expected) {
//				t.Errorf("expected %v, got %v", tt.expected, got)
//			}
//		})
//	}
//}
//
//func TestFindByWeek(t *testing.T) {
//	storage := New()
//	// 2023-10-02 12:00:00 +0000 UTC (Monday)
//	start := time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC)
//	// 2023-10-02 14:00:00 +0000 UTC
//	end := start.Add(2 * time.Hour)
//	storage.Add(StartEndDt{StartDt: start, EndDt: end})
//
//	tests := []struct {
//		name     string
//		date     time.Time
//		expected map[int]interface{}
//	}{
//		{
//			name: "Find event in the same week",
//			date: start,
//			expected: map[int]interface{}{
//				1: StartEndDt{StartDt: start, EndDt: end},
//			},
//		},
//		{
//			name:     "Find event in a different week",
//			date:     start.AddDate(0, 0, 7),
//			expected: map[int]interface{}{},
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//			got := storage.ListByWeek(tt.date)
//			if !equal(got, tt.expected) {
//				t.Errorf("expected %v, got %v", tt.expected, got)
//			}
//		})
//	}
//}
//
//func TestFindByMonth(t *testing.T) {
//	storage := New()
//	// 2023-10-01 12:00:00 +0000 UTC
//	start := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
//	// 2023-10-01 14:00:00 +0000 UTC
//	end := start.Add(2 * time.Hour)
//	storage.Add(context.Background(), storage2.Event{StartDt: start, EndDt: end})
//
//	tests := []struct {
//		name     string
//		date     time.Time
//		expected map[int]interface{}
//	}{
//		{
//			name: "Find event in the same month",
//			date: start,
//			expected: map[int]interface{}{
//				1: StartEndDt{StartDt: start, EndDt: end},
//			},
//		},
//		{
//			name:     "Find event in a different month",
//			date:     start.AddDate(0, 1, 0),
//			expected: map[int]interface{}{},
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//			got := storage.ListByMonth(ctx, tt.date)
//			if !equal(got, tt.expected) {
//				t.Errorf("expected %v, got %v", tt.expected, got)
//			}
//		})
//	}
//}

func createStubEvent(name string, startDt time.Time, endDt time.Time) storage.Event {
	uuid, _ := uuid.NewV7()
	if startDt.IsZero() {
		startDt = time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	}
	if endDt.IsZero() {
		endDt = startDt.Add(time.Hour)
	}
	return storage.Event{
		ID:          uuid,
		Title:       name,
		StartDt:     startDt,
		EndDt:       endDt,
		Description: "description",
		UserId:      1,
		Notify:      time.Second,
	}
}
