package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	logg       Logger
	app        Application
	httpServer *http.Server
}

type Logger interface {
	InfoRaw(msg string)
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logg Logger, app Application, conf config.RestServerConf) *Server {
	return &Server{
		logg: logg,
		app:  app,
		httpServer: &http.Server{
			Addr:        fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			ReadTimeout: 10 * time.Second,
			IdleTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("Starting HTTP server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	loggedMux := loggingMiddleware(mux, s.logg)

	s.httpServer.Handler = loggedMux

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logg.Error("HTTP server ListenAndServe: " + err.Error())
		}
	}()

	<-ctx.Done()
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stopping HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logg.Error("HTTP server Shutdown: " + err.Error())
		return err
	}

	s.logg.Info("HTTP server stopped")
	return nil
}
