package memorystorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model/event"
	"github.com/google/uuid"
)

var ctx = context.Background()

func equal(a, b []event.Event) bool {
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
	e := createStubEvent("event 1", time.Time{}, time.Time{})

	tests := []struct {
		name      string
		operation func(s *Storage) error
		expected  []event.Event
	}{
		{
			name: "Add event",
			operation: func(s *Storage) error {
				s.Add(ctx, e)
				return nil
			},
			expected: []event.Event{e},
		},
		{
			name: "Delete event",
			operation: func(s *Storage) error {
				s.Add(ctx, e)
				return s.Delete(ctx, e.ID)
			},
			expected: []event.Event{},
		},
		{
			name: "Delete non-existent event",
			operation: func(s *Storage) error {
				id, _ := uuid.NewV7()
				return s.Delete(ctx, id)
			},
			expected: []event.Event{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := New()
			err := tt.operation(memStorage)
			if err != nil && !errors.Is(err, adapters.ErrNotFound) {
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
	e := createStubEvent("event 1", time.Time{}, time.Time{})
	memStorage.Add(ctx, e)

	e.Title = "event 2"
	_, err := memStorage.Update(ctx, e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateNotExistentEvent(t *testing.T) {
	memStorage := New()
	e := createStubEvent("event 1", time.Time{}, time.Time{})
	_, err := memStorage.Update(ctx, e)
	if !errors.Is(err, adapters.ErrNotFound) {
		t.Errorf("expected error %v, got %v", adapters.ErrNotFound, err)
	}
}

func TestUpdate(t *testing.T) {
	memStorage := New()
	// 2000-01-01 12:00:00 +0000 UTC
	start := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	// 2000-01-01 14:00:00 +0000 UTC
	end := start.Add(2 * time.Hour)
	e1 := createStubEvent("event 1", start, end)

	e2 := e1
	e2.Title = "event 2"

	tests := []struct {
		name      string
		operation func() error
		expected  []event.Event
	}{
		{
			name: "Update item",
			operation: func() error {
				memStorage = New()
				memStorage.Add(ctx, e1)
				_, err := memStorage.Update(ctx, e2)
				return err
			},
			expected: []event.Event{e2},
		},
		{
			name: "Update non-existent item",
			operation: func() error {
				memStorage = New()
				event := createStubEvent("not exist event", start, end)
				_, err := memStorage.Update(ctx, event)
				return err
			},
			expected: []event.Event{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if err != nil && !errors.Is(err, adapters.ErrNotFound) {
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

	e1 := createStubEvent("event 1", start, end)
	memStorage.Add(ctx, e1)

	tests := []struct {
		name     string
		date     time.Time
		expected []event.Event
	}{
		{
			name:     "List existing event",
			date:     start,
			expected: []event.Event{e1},
		},
		{
			name:     "List non-existing event",
			date:     start.Add(24 * time.Hour),
			expected: []event.Event{},
		},
	}

	for _, tt := range tests {
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

	e := createStubEvent("event 1", start, end)

	memStorage.Add(ctx, e)

	tests := []struct {
		name     string
		startDt  time.Time
		endDt    time.Time
		expected []event.Event
	}{
		{
			name:     "List overlapping event",
			startDt:  start.Add(-1 * time.Hour),
			endDt:    start.Add(1 * time.Hour),
			expected: []event.Event{e},
		},
		{
			name:     "List non-overlapping event",
			startDt:  start.Add(3 * time.Hour),
			endDt:    start.Add(4 * time.Hour),
			expected: []event.Event{},
		},
		{
			name:     "List event within range",
			startDt:  start,
			endDt:    end,
			expected: []event.Event{e},
		},
	}

	for _, tt := range tests {
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
	e := createStubEvent("event 1", start, end)
	memStorage.Add(ctx, e)

	tests := []struct {
		name     string
		date     time.Time
		expected []event.Event
	}{
		{
			name:     "List events in the same week",
			date:     start,
			expected: []event.Event{e},
		},
		{
			name:     "List events in a different week",
			date:     start.AddDate(0, 0, 7),
			expected: []event.Event{},
		},
	}

	for _, tt := range tests {
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
	e := createStubEvent("event 1", start, end)
	memStorage.Add(ctx, e)

	tests := []struct {
		name     string
		date     time.Time
		expected []event.Event
	}{
		{
			name:     "List event in the same month",
			date:     start,
			expected: []event.Event{e},
		},
		{
			name:     "List event in a different month",
			date:     start.AddDate(0, 1, 0),
			expected: []event.Event{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := memStorage.ListByMonth(ctx, tt.date)
			if !equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func createStubEvent(name string, startDt time.Time, endDt time.Time) event.Event {
	eventID, _ := uuid.NewV7()
	userID, _ := uuid.NewV7()
	if startDt.IsZero() {
		startDt = time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	}
	if endDt.IsZero() {
		endDt = startDt.Add(time.Hour)
	}
	return event.Event{
		ID:          eventID,
		Title:       name,
		StartDt:     startDt,
		EndDt:       endDt,
		Description: "description",
		UserID:      userID,
		Notify:      time.Second,
	}
}
