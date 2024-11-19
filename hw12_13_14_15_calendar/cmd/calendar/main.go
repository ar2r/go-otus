package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	memorystorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/memory"
	sqlstorage "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/adapters/pgx"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model/event"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
	pb "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/protobuf"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/services"
)

var (
	configFile string
	logg       *logger.Logger
)

type pb_server struct {
	pb.UnimplementedEventServiceServer
}

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"/etc/calendar/config.toml.example",
		"Path to configuration file",
	)
}

func main() {
	mainCtx := context.Background()
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
		eventRepo, err = sqlstorage.New(mainCtx, myConfig.Database, logg)
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

	// Signal handler
	mainCtx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	logg.Info("Signal handler initialized")

	wg := sync.WaitGroup{}
	wg.Add(2)

	// REST server
	server := internalhttp.NewServer(logg, calendar, myConfig.RestServer)
	logg.Info("HTTP server initialized")

	go func() {
		if err := server.Start(mainCtx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
		wg.Done()
	}()

	// GRPC server
	service := services.NewEventService(eventRepo)
	grpcServerService := grpc.NewService(service)
	grpcServer := grpc.NewServer(myConfig.GrpcServer, grpcServerService)
	logg.Info("GRPC server initialized")

	go func() {
		if err := grpcServer.Run(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
		wg.Done()
	}()

	// Graceful shutdown
	go func() {
		<-mainCtx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("Calendar is running...")
	wg.Wait()
	logg.Info("App shutdown!")
}

func InitLogger(loggerConf config.LoggerConf) *logger.Logger {
	return logger.New(
		loggerConf.Level,
		loggerConf.Channel,
		loggerConf.Filename,
	)
}
