package appSender

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue"
)

type AppSender struct {
	logg       *slog.Logger
	conf       *config.Config
	repository model.EventRepository
	consumer   queue.IConsumer
}

func NewSender(logg *slog.Logger, conf *config.Config, repo model.EventRepository) *AppSender {
	return &AppSender{
		logg:       logg,
		conf:       conf,
		repository: repo,
	}
}

func (a *AppSender) Run() error {
	eventsCh := make(chan string)
	doneCh := make(chan error)

	queueConn, err := queue.NewConsumer(a.logg, a.conf, doneCh)
	if err != nil {
		return fmt.Errorf("failed to create queue consumer: %w", err)
	}
	a.consumer = queueConn

	go handle(a.logg, eventsCh, doneCh)

	err = queueConn.Consume(eventsCh)
	if err != nil {
		return fmt.Errorf("failed to consume from queue: %w", err)
	}

	a.logg.Info("Sender is running...")
	return nil
}

func handle(logg *slog.Logger, messages <-chan string, done chan<- error) {
	for m := range messages {
		event := model.Event{}
		json.Unmarshal([]byte(m), &event)
		logg.Warn("You have got notified about event: " + event.Title)
	}
	logg.Debug("handle: deliveries channel closed")
	done <- nil
}
