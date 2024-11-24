package appscheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/go-co-op/gocron/v2"
)

type AppScheduler struct {
	logg      *slog.Logger
	conf      *config.Config
	repo      model.EventRepository
	producer  queue.IProducer
	scheduler gocron.Scheduler
}

func New(logg *slog.Logger, conf *config.Config, repo model.EventRepository, producer queue.IProducer) (*AppScheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &AppScheduler{
		logg:      logg,
		conf:      conf,
		repo:      repo,
		producer:  producer,
		scheduler: scheduler,
	}, nil
}

func (a *AppScheduler) Run(ctx context.Context) error {
	task := gocron.NewTask(
		a.produceNotification(ctx),
		a.producer,
		a.conf.RabbitMQ.RoutingKey,
	)

	_, err := a.scheduler.NewJob(
		gocron.CronJob("1/2 * * * * *", true),
		task,
		gocron.WithName("notify"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		a.logg.Error("Failed to create job: " + err.Error())
	}

	a.scheduler.Start()
	a.logg.Info("Scheduler is running...")
	return nil
}

func (a *AppScheduler) produceNotification(ctx context.Context) func(queueConn queue.IProducer, routingKey string) {
	return func(queueConn queue.IProducer, routingKey string) {
		events, err := a.repo.ListNotNotified(ctx)
		if err != nil {
			a.logg.Error("failed to get events: " + err.Error())
			return
		}
		for _, event := range events {
			eventJSON, err := json.Marshal(event)
			if err != nil {
				a.logg.Error("failed to marshal event: " + err.Error())
				continue
			}
			if err = queueConn.Publish(routingKey, eventJSON); err != nil {
				a.logg.Error("failed to publish message: " + err.Error())
			}
			event.NotificationSent = true

			_, err = a.repo.Update(ctx, event)
			if err != nil {
				a.logg.Error("failed to update event: " + err.Error())
			}
			a.logg.Debug("Event has been sent to MQ: " + event.ID.String())
		}
	}
}

func (a *AppScheduler) Stop() error {
	err := a.scheduler.StopJobs()
	if err != nil {
		return fmt.Errorf("failed to stop jobs: %w", err)
	}
	return nil
}
