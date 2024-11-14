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

	event := createStubEvent("event 1", time.Time{}, time.Time{})

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
	event := createStubEvent("event 1", time.Time{}, time.Time{})
	memStorage.Add(ctx, event)

	event.Title = "item2"
	_, err := memStorage.Update(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateNotExistentEvent(t *testing.T) {
	memStorage := New()
	event := createStubEvent("item 1", time.Time{}, time.Time{})
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
	event1 := createStubEvent("event 1", start, end)

	event2 := event1
	event2.Title = "event 2"

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

func TestListByDate(t *testing.T) {
	memStorage := New()

	// 2000-01-01 12:00:00 +0000 UTC
	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)

	event1 := createStubEvent("event 1", start, end)
	memStorage.Add(ctx, event1)

	tests := []struct {
		name     string
		date     time.Time
		expected []storage.Event
	}{
		{
			name:     "List existing event",
			date:     start,
			expected: []storage.Event{event1},
		},
		{
			name:     "List non-existing event",
			date:     start.Add(24 * time.Hour),
			expected: []storage.Event{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, _ := memStorage.ListByDate(ctx, tt.date)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestListByPeriod(t *testing.T) {
	memStorage := New()
	// 2000-01-01 12:00:00 +0000 UTC
	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)

	event := createStubEvent("event 1", start, end)

	memStorage.Add(ctx, event)

	tests := []struct {
		name     string
		startDt  time.Time
		endDt    time.Time
		expected []storage.Event
	}{
		{
			name:     "List overlapping event",
			startDt:  start.Add(-1 * time.Hour),
			endDt:    start.Add(1 * time.Hour),
			expected: []storage.Event{event},
		},
		{
			name:     "List non-overlapping event",
			startDt:  start.Add(3 * time.Hour),
			endDt:    start.Add(4 * time.Hour),
			expected: []storage.Event{},
		},
		{
			name:     "List event within range",
			startDt:  start,
			endDt:    end,
			expected: []storage.Event{event},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, _ := memStorage.ListByPeriod(ctx, tt.startDt, tt.endDt)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestListByWeek(t *testing.T) {
	memStorage := New()
	// 2023-10-02 12:00:00 +0000 UTC (Monday)
	start := time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC)
	// 2023-10-02 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)
	event := createStubEvent("event 1", start, end)
	memStorage.Add(ctx, event)

	tests := []struct {
		name     string
		date     time.Time
		expected []storage.Event
	}{
		{
			name:     "List events in the same week",
			date:     start,
			expected: []storage.Event{event},
		},
		{
			name:     "List events in a different week",
			date:     start.AddDate(0, 0, 7),
			expected: []storage.Event{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, _ := memStorage.ListByWeek(ctx, tt.date)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestListByMonth(t *testing.T) {
	memStorage := New()
	// 2023-10-01 12:00:00 +0000 UTC
	start := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	// 2023-10-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)
	event := createStubEvent("event 1", start, end)
	memStorage.Add(ctx, event)

	tests := []struct {
		name     string
		date     time.Time
		expected []storage.Event
	}{
		{
			name:     "List event in the same month",
			date:     start,
			expected: []storage.Event{event},
		},
		{
			name:     "List event in a different month",
			date:     start.AddDate(0, 1, 0),
			expected: []storage.Event{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, _ := memStorage.ListByMonth(ctx, tt.date)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func createStubEvent(name string, startDt time.Time, endDt time.Time) storage.Event {
	eventId, _ := uuid.NewV7()
	userId, _ := uuid.NewV7()
	if startDt.IsZero() {
		startDt = time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	}
	if endDt.IsZero() {
		endDt = startDt.Add(time.Hour)
	}
	return storage.Event{
		ID:          eventId,
		Title:       name,
		StartDt:     startDt,
		EndDt:       endDt,
		Description: "description",
		UserId:      userId,
		Notify:      time.Second,
	}
}
