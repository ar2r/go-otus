package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/sender"
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

	// Config
	myConfig, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Logger
	logg = myslog.New(
		myConfig.Logger.Level,
		myConfig.Logger.Channel,
		myConfig.Logger.Filename,
	)

	logg.Info("Sender is booting...")

	// Signal handler
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// App
	app := sender.New(logg, myConfig, eventRepo)
	err = app.Run()
	if err != nil {
		logg.Error("failed to run sender: " + err.Error())
		return
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		err = app.Stop()
		if err != nil {
			logg.Error("failed to stop consumer: " + err.Error())
		}
	}()

	<-ctx.Done()
	logg.Info("Sender shutdown!")
}
