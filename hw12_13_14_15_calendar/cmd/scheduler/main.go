package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/appScheduler"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/myslog"
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

	// Queue producer
	producerConn, err := queue.NewProducer(logg, myConfig)
	if err != nil {
		logg.Error("failed to create queue producer: " + err.Error())
		return
	}

	// App
	app, err := appScheduler.New(logg, myConfig, eventRepo, producerConn)
	if err != nil {
		logg.Error("failed to create app: " + err.Error())
		return
	}

	if err = app.Run(ctx); err != nil {
		logg.Error("failed to run app: " + err.Error())
		return
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		err = app.Stop()
		if err != nil {
			logg.Error("failed to stop cron jobs: " + err.Error())
		}
	}()

	<-ctx.Done()
	logg.Info("Scheduler shutdown!")
}
