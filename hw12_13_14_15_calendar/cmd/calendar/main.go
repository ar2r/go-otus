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
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/db"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"
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
	var storage app.Storage

	ctx := context.Background()
	flag.Parse()

	if flag.Arg(0) == "version" {
		VersionPrint()
		return
	}

	myConfig, err := InitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}

	if flag.Arg(0) == "migrate" {
		if err := MigrateRun(logg, myConfig.Database, true); err != nil {
			logg.Error(fmt.Sprintf("%v", err))
		}
		return
	}

	// UserRepository
	switch myConfig.App.Storage {
	case "memory":
		storage = memorystorage.New()
		logg.Info("Memory adapters initialized")
	case "pgx":
		pgxPool, err := db.Connect(ctx, myConfig.Database, logg)
		defer func() {
			if pgxPool != nil {
				db.Close(pgxPool)
			}
		}()
		if err != nil {
			logg.Error(fmt.Sprintf("failed to create connetion to dictionaries db: %s", err))
		}
		storage = sqlstorage.New(pgxPool)
		logg.Info("SQL adapters initialized")
	default:
		logg.Error("Invalid adapters type: " + myConfig.App.Storage)
		return
	}

	// Application
	calendar := app.New(logg, storage)
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

func InitLogger() (*config.Config, error) {
	// Обработка конфига
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logg = logger.New(
		conf.Logger.Level,
		conf.Logger.Channel,
		conf.Logger.Filename,
	)

	if conf.App.Debug {
		logg.Report()
	}
	return conf, err
}
