package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	memorystorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/memory"
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/pgx"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model/event"
	internalhttp "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
)

var (
	configFile string
	logg       *logger.Logger
)

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"/etc/calendar/config.toml.example",
		"Path to configuration file",
	)
}

func main() {
	ctx := context.Background()
	flag.Parse()

	if flag.Arg(0) == "version" {
		VersionPrint()
		return
	}

	myConfig, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logg = InitLogger(myConfig.Logger)
	if myConfig.App.Debug {
		logg.Report()
	}

	if flag.Arg(0) == "migrate" {
		if err := MigrateRun(logg, myConfig.Database, true); err != nil {
			logg.Error(fmt.Sprintf("%v", err))
		}
		return
	}

	// Event Repository
	var eventRepo event.Repository

	switch myConfig.App.Storage {
	case "memory":
		eventRepo = memorystorage.New()
		logg.Info("Memory adapters initialized")
	case "sql":
		eventRepo, err = sqlstorage.New(ctx, myConfig.Database, logg)
		if err != nil {
			logg.Error(fmt.Sprintf("failed to initialize SQL storage: %s", err))
			return
		}
		logg.Info("SQL adapters initialized")
	default:
		logg.Error("Invalid adapters type: " + myConfig.App.Storage)
		return
	}

	// Application
	calendar := app.New(
		logg,
		eventRepo,
	)
	logg.Info("App initialized")

	// HTTP server
	server := internalhttp.NewServer(logg, calendar, myConfig.Server)
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
		return
	}
}

func InitLogger(loggerConf config.LoggerConf) *logger.Logger {
	return logger.New(
		loggerConf.Level,
		loggerConf.Channel,
		loggerConf.Filename,
	)
}
