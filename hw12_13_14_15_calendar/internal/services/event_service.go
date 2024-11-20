package services

import (
	"context"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/dto"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
)

type EventServiceInterface interface {
	Add(ctx context.Context, dto dto.CreateEventDto) (*model.Event, error)
	Update(ctx context.Context, dto dto.UpdateEventDto) (*model.Event, error)
	Delete(ctx context.Context, dto dto.DeleteEventDto) error
	ListByDate(ctx context.Context, dto dto.ListByDateDto) ([]model.Event, error)
	ListByWeek(ctx context.Context, dto dto.ListByDateDto) ([]model.Event, error)
	ListByMonth(ctx context.Context, dto dto.ListByDateDto) ([]model.Event, error)
}

type EventService struct {
	repository model.EventRepository
}

func NewEventService(repository model.EventRepository) EventServiceInterface {
	return &EventService{
		repository: repository,
	}
}

func (s EventService) Add(ctx context.Context, dto dto.CreateEventDto) (*model.Event, error) {
	e := dto.ToModel()
	err := e.GenerateID()
	if err != nil {
		return nil, err
	}
	return s.repository.Add(ctx, e)
}

func (s EventService) Update(ctx context.Context, dto dto.UpdateEventDto) (*model.Event, error) {
	um := dto.ToModel()
	return s.repository.Update(ctx, um)
}

func (s EventService) Delete(ctx context.Context, dto dto.DeleteEventDto) error {
	return s.repository.Delete(ctx, dto.ID)
}

func (s EventService) ListByDate(ctx context.Context, dto dto.ListByDateDto) ([]model.Event, error) {
	return s.repository.ListByDate(ctx, dto.Date)
}

func (s EventService) ListByWeek(ctx context.Context, dto dto.ListByDateDto) ([]model.Event, error) {
	return s.repository.ListByWeek(ctx, dto.Date)
}

func (s EventService) ListByMonth(ctx context.Context, dto dto.ListByDateDto) ([]model.Event, error) {
	return s.repository.ListByMonth(ctx, dto.Date)
}
