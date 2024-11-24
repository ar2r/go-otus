package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	// Signal handler
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	logg.Info("Signal handler initialized")

	// Queue producer
	eventsCh := make(chan string)
	doneCh := make(chan error)

	queueConn, err := queue.NewConsumer(logg, myConfig, doneCh)
	if err != nil {
		logg.Error("failed to create queue consumer: " + err.Error())
		return
	}

	go handle(logg, eventsCh, doneCh)

	err = queueConn.Consume(eventsCh)
	if err != nil {
		logg.Error("failed to consume from queue: " + err.Error())
		return
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()

		// Разорвать коннект с кроликом
	}()
	logg.Info("Sender is running...")

	// Wait for signal
	<-ctx.Done()

	logg.Info("Sender shutdown!")
}

// todo: Перенести
func handle(logg *slog.Logger, messages <-chan string, done chan<- error) {
	for m := range messages {
		logg.Warn("got message: " + m)
	}
	logg.Info("handle: deliveries channel closed")
	done <- nil
}
