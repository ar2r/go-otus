package appScheduler

import (
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue"
)

type AppScheduler struct {
	repo     model.EventRepository
	producer queue.IProducer
}

func NewScheduler(repo model.EventRepository, producer queue.IProducer) *AppScheduler {
	return &AppScheduler{
		repo:     repo,
		producer: producer,
	}
}

func (a *AppScheduler) Run() error {
	// Тут читать из БД события, которые нужно отправить

	return nil
}
