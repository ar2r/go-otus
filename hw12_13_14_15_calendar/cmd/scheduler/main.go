package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage"
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
	errorCh := make(chan error) // Канал для передачи ошибок между горутинами

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

	logg.Info("Scheduler is booting...")

	// Event EventRepository
	eventRepo, err = storage.New(ctx, logg, myConfig)
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
	app, err := scheduler.New(logg, myConfig, eventRepo, producerConn, errorCh)
	if err != nil {
		logg.Error("failed to create app: " + err.Error())
		return
	}

	if err = app.Run(ctx); err != nil {
		logg.Error("failed to run app: " + err.Error())
		return
	}

	go func() {
		for {
			select {
			case <-errorCh:
				cancel()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		producerConn.Close()
		app.Stop()
	}()

	<-ctx.Done()
	logg.Info("Scheduler shutdown!")
}
