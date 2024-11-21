package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
)

// Server HTTP сервер для обработки REST запросов.
type Server struct {
	app        app.IApplication
	logg       IServerLogger
	httpServer *http.Server
}

type IServerLogger interface {
	Info(msg string, attrs ...interface{})
	Error(msg string, attrs ...interface{})
}

func NewServer(app app.IApplication, logg IServerLogger, conf Config) *Server {
	return &Server{
		app:  app,
		logg: logg,
		httpServer: &http.Server{
			Addr:        fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			ReadTimeout: 10 * time.Second,
			IdleTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("Starting HTTP server...")
	s.registerLogger(s.registerRoutes())

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logg.Error("HTTP server ListenAndServe: " + err.Error())
		}
	}()

	s.logg.Info("HTTP Waiting ctx done")
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

func (s *Server) registerRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("POST /events", s.createEventHandler)
	mux.HandleFunc("PUT /events/{id}", s.updateEventHandler)
	mux.HandleFunc("DELETE /events/{id}", s.deleteEventHandler)
	return mux
}

func (s *Server) registerLogger(mux *http.ServeMux) {
	s.httpServer.Handler = loggingMiddleware(mux, s.logg)
}
