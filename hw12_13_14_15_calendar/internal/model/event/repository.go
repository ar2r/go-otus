package event

import (
	"context"
	"time"

	model2 "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	CRUDer
	Lister
}

type CRUDer interface {
	Add(ctx context.Context, e model2.Event) (*model2.Event, error)
	Update(ctx context.Context, e model2.Event) (*model2.Event, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	Get(ctx context.Context, uuid uuid.UUID) (*model2.Event, error)
}

type Lister interface {
	List(ctx context.Context) ([]model2.Event, error)
	ListByDate(ctx context.Context, start time.Time) ([]model2.Event, error)
	ListByPeriod(ctx context.Context, startDt time.Time, endDt time.Time) ([]model2.Event, error)
	ListByWeek(ctx context.Context, startDt time.Time) ([]model2.Event, error)
	ListByMonth(ctx context.Context, startDt time.Time) ([]model2.Event, error)
}
