package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Обработка конфига
	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Logger
	logg := logger.New(
		config.Logger.Level,
		config.Logger.Channel,
		config.Logger.Filename,
	)

	if config.App.Debug {
		logg.Report()
	}

	// Storage
	storage := memorystorage.New()
	logg.Warn("No storage initialized")

	// Application
	calendar := app.New(logg, storage)
	logg.Info("App initialized")

	// HTTP server
	server := internalhttp.NewServer(logg, calendar)
	logg.Info("HTTP server initialized")

	// Signal handler
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	logg.Info("Signal handler initialized")

	// Graceful shutdown
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("Calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
