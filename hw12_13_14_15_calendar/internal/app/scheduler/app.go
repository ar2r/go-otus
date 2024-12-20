package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/go-co-op/gocron/v2"
)

const (
	CrontabNotify = "1/2 * * * * *"
	CrontabClean  = "1 * * * * *"
)

type AppScheduler struct {
	logg      *slog.Logger
	conf      *config.Config
	repo      model.EventRepository
	producer  MessageProducer
	scheduler gocron.Scheduler
	errorCh   chan<- error
}

type MessageProducer interface {
	Publish(routingKey string, body []byte) error
	Close()
}

func New(logg *slog.Logger, conf *config.Config, repo model.EventRepository, producer MessageProducer, errorCh chan<- error) (*AppScheduler, error) {
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
		errorCh:   errorCh,
	}, nil
}

func (a *AppScheduler) Run(ctx context.Context) error {
	a.registerProducerTask(ctx)
	a.registerCleanerTask(ctx)

	a.scheduler.Start()
	a.logg.Info("Scheduler is running...")
	return nil
}

func (a *AppScheduler) Stop() {
	err := a.scheduler.StopJobs()
	if err != nil {
		a.logg.Warn("Failed to stop jobs: " + err.Error())
	}
}

func (a *AppScheduler) registerProducerTask(ctx context.Context) {
	produceTask := gocron.NewTask(
		a.produceNotification(ctx),
		a.producer,
		a.conf.RabbitMQ.RoutingKey,
	)
	_, err := a.scheduler.NewJob(
		gocron.CronJob(CrontabNotify, true),
		produceTask,
		gocron.WithName("notify"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		a.logg.Error("Failed to create job: " + err.Error())
	}
}

func (a *AppScheduler) registerCleanerTask(ctx context.Context) {
	cleanerTask := gocron.NewTask(
		a.cleanupEvents(ctx),
		a.conf.App.CleanupDuration,
	)
	_, err := a.scheduler.NewJob(
		gocron.CronJob(CrontabClean, true),
		cleanerTask,
		gocron.WithName("cleaner"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		a.logg.Error("Failed to create job: " + err.Error())
	}
}
