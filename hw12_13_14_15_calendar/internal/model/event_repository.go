package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventRepository interface {
	cruder
	lister
	cleaner
}

type cruder interface {
	Add(ctx context.Context, e Event) (Event, error)
	Update(ctx context.Context, e Event) (Event, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	Get(ctx context.Context, uuid uuid.UUID) (Event, error)
}

type lister interface {
	ListByDate(ctx context.Context, start time.Time) ([]Event, error)
	ListByPeriod(ctx context.Context, startDt time.Time, endDt time.Time) ([]Event, error)
	ListByWeek(ctx context.Context, startDt time.Time) ([]Event, error)
	ListByMonth(ctx context.Context, startDt time.Time) ([]Event, error)
	ListNotNotified(ctx context.Context) ([]Event, error)
}

type cleaner interface {
	DeleteOlderThan(ctx context.Context, t time.Time) error
}
