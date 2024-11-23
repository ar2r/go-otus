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

	// Graceful shutdown
	go func() {
		<-ctx.Done()

		// Разорвать коннект с кроликом
	}()
	logg.Info("Sender is running...")

	// todo: Тут читать из кролика и выводить сообщения в консоль )))
	// todo: Запустить горутинку, которая будет читать из кролика и выводить сообщения в консоль

	// Wait for signal
	<-ctx.Done()

	logg.Info("Sender shutdown!")
}
