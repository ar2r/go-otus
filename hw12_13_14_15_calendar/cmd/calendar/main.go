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

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
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

	if flag.Arg(0) == "version" {
		VersionPrint()
		return
	}

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

	logg.Info("Calendar is booting...")

	if flag.Arg(0) == "migrate" {
		if err := MigrateRun(logg, myConfig.Database, true); err != nil {
			logg.Error(fmt.Sprintf("%v", err))
		}
		return
	}

	// Event EventRepository
	eventRepo, err = adapters.New(ctx, logg, myConfig)
	if err != nil {
		logg.Error("failed to create repository: " + err.Error())
		return
	}

	// Application
	calendar := calendar.New(eventRepo)
	logg.Info("App initialized")

	// Signal handler
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	logg.Info("Signal handler initialized")

	// Servers
	serversWG := sync.WaitGroup{}
	serversWG.Add(2)

	// REST server
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

	// GRPC server
	grpcServerService := grpcserver.NewService(calendar)
	grpcServer := grpcserver.NewServer(logg, myConfig.GRPCServer, grpcServerService)
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

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop HTTP server: " + err.Error())
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop GRPC server: " + err.Error())
		}
	}()
	logg.Info("Calendar is running...")

	serversWG.Wait()
	logg.Info("App shutdown!")
}
