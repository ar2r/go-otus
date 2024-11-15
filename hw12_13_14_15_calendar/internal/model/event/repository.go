package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const CtxKeyUserID = "user_id"

type Interface interface {
	CRUDer
	Creator
	Lister
}
type Creator interface {
	CreateEvent(ctx context.Context, id uuid.UUID, title string) error
}

type CRUDer interface {
	Add(ctx context.Context, e Event) (*Event, error)
	Update(ctx context.Context, e Event) (*Event, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	Get(ctx context.Context, uuid uuid.UUID) (*Event, error)
}

type Lister interface {
	List(ctx context.Context) ([]Event, error)
	ListByDate(ctx context.Context, start time.Time) ([]Event, error)
	ListByPeriod(ctx context.Context, startDt time.Time, endDt time.Time) ([]Event, error)
	ListByWeek(ctx context.Context, startDt time.Time) ([]Event, error)
	ListByMonth(ctx context.Context, startDt time.Time) ([]Event, error)
}
