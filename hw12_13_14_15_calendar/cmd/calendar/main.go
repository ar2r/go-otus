package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	memorystorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/memory"
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/pgx"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/services"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/pkg/easylog"
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
		"/etc/calendar/config.toml.example",
		"Path to configuration file",
	)
}

func GetLogger() *slog.Logger {
	return logg
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

	logg = initLogger(myConfig.Logger)

	if flag.Arg(0) == "migrate" {
		if err := MigrateRun(logg, myConfig.Database, true); err != nil {
			logg.Error(fmt.Sprintf("%v", err))
		}
		return
	}

	// Event EventRepository
	eventRepo, err = initRepository(ctx, myConfig)
	if err != nil {
		logg.Error("failed to create repository: " + err.Error())
		return
	}

	// Application
	calendar := app.New(logg, eventRepo)
	logg.Info("App initialized")

	// Signal handler
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	logg.Info("Signal handler initialized")

	/**
	 *
	 * Start httpServer listeners
	 *
	 */
	serversWG := sync.WaitGroup{}
	serversWG.Add(2)

	// REST httpServer
	httpServer := internalhttp.NewServer(calendar, logg, myConfig.HTTPServer)
	logg.Info("HTTP server initialized")

	go func() {
		defer serversWG.Done()

		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start HTTP server: " + err.Error())
			cancel()
			return
		}
	}()

	// GRPC httpServer
	service := services.NewEventService(eventRepo)
	grpcServerService := grpcserver.NewService(service)
	grpcServer := grpcserver.NewServer(myConfig.GRPCServer, grpcServerService)
	logg.Info("GRPC server initialized")

	go func() {
		defer serversWG.Done()

		if err := grpcServer.Run(); err != nil {
			logg.Error("failed to start grpc GRPC server: " + err.Error())
			cancel()
		}
	}()

	// Graceful shutdown
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop HTTP server: " + err.Error())
		}
		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop GRPC server: " + err.Error())
		}
	}()
	logg.Info("Calendar is running...")

	serversWG.Wait()
	logg.Info("App shutdown!")
}

func initRepository(ctx context.Context, myConfig *config.Config) (model.EventRepository, error) {
	var eventRepo model.EventRepository
	var err error

	switch myConfig.App.Storage {
	case "memory":
		eventRepo = memorystorage.New()
		logg.Info("Memory adapters initialized")
	case "sql":
		eventRepo, err = sqlstorage.New(ctx, myConfig.Database, logg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SQL storage: %w", err)
		}
		logg.Info("SQL adapters initialized")
	default:
		return nil, fmt.Errorf("invalid adapters type: %s", myConfig.App.Storage)
	}
	return eventRepo, nil
}

func initLogger(loggerConf config.LoggerConfig) *slog.Logger {
	return easylog.New(
		loggerConf.Level,
		loggerConf.Channel,
		loggerConf.Filename,
	)
}
