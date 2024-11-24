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

func New(logg *slog.Logger, conf *config.Config, repo model.EventRepository) *AppSender {
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

func (a *AppSender) Stop() error {
	err := a.consumer.Close()
	if err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}
	return nil
}

// handle Вывести в консоль уведомление о предстоящих событиях.
func handle(logg *slog.Logger, messages <-chan string, done chan<- error) {
	for m := range messages {
		event := model.Event{}
		err := json.Unmarshal([]byte(m), &event)
		if err != nil {
			logg.Error("failed to unmarshal event: " + err.Error())
			continue
		}
		logg.Warn("Event starts at: " + event.StartDt.String() + " title: " + event.Title)
	}
	done <- nil
}
