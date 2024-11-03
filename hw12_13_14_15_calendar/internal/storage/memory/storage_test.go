package memorystorage

import (
	"errors"
	"testing"
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
			name: "Delete item",
			operation: func(s *Storage) error {
				id := s.Add("item1")
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
