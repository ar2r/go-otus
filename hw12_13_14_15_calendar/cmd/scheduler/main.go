package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/myslog"
	"github.com/go-co-op/gocron/v2"
)

var (
	configFile string
	logg       *slog.Logger
	eventRepo  model.EventRepository
)

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"/etc/calendar/config.toml",
		"Path to configuration file",
	)
}

func main() {
	ctx := context.Background()
	flag.Parse()

	myConfig, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logg = myslog.New(
		myConfig.Logger.Level,
		myConfig.Logger.Channel,
		myConfig.Logger.Filename,
	)

	// Event EventRepository
	eventRepo, err = adapters.New(ctx, logg, myConfig)
	if err != nil {
		logg.Error("failed to create repository: " + err.Error())
		return
	}

	// Signal handler
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	logg.Info("Signal handler initialized")

	s, err := gocron.NewScheduler()
	if err != nil {
		logg.Error("Failed to create scheduler: ", err)
		os.Exit(1)
	}

	task := gocron.NewTask(
		func(a string, b int) {
			// do things
			logg.Debug("Task executed1: ", a, b)
			fmt.Printf("Task executed2: %s %d\n", a, b)
			time.Sleep(5 * time.Second)
		},
		"hello",
		1,
	)

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.CronJob("1/2 * * * * *", true),
		task,
		gocron.WithName("find-notify-events"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		logg.Error("Failed to create job: ", err)
	}

	s.Start()

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		err = s.StopJobs()
		if err != nil {
			logg.Error("Failed to stop jobs: ", err)
		}
	}()
	logg.Info("Scheduler is running...")

	// Wait for signal
	<-ctx.Done()

	logg.Info("Scheduler shutdown!")
}
