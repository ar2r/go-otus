package dto

import (
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type CreateEventDto struct {
	UserID      uuid.UUID     `json:"userId"`
	Title       string        `json:"title"`
	StartDt     time.Time     `json:"startDt"`
	EndDt       time.Time     `json:"endDt"`
	Description string        `json:"description"`
	NotifyAt    time.Duration `json:"notifyAt"`
}

func (d *CreateEventDto) ToModel() model.Event {
	return model.Event{
		UserID:      d.UserID,
		Title:       d.Title,
		StartDt:     d.StartDt,
		EndDt:       d.EndDt,
		Description: d.Description,
		NotifyAt:    d.NotifyAt,
	}
}
