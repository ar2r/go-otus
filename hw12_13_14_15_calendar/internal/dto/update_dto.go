package dto

import (
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

type UpdateEventDto struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.UUID     `json:"userId"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	StartDt     time.Time     `json:"startDt"`
	EndDt       time.Time     `json:"endDt"`
	NotifyAt    time.Duration `json:"notifyAt"`
}

func (u UpdateEventDto) ToModel() model.Event {
	return model.Event{
		ID:          u.ID,
		UserID:      u.UserID,
		Title:       u.Title,
		Description: u.Description,
		StartDt:     u.StartDt,
		EndDt:       u.EndDt,
		NotifyAt:    u.NotifyAt,
	}
}
