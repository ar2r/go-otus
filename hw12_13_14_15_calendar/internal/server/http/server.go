package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
)

// Server HTTP сервер для обработки REST запросов.
type Server struct {
	app        *app.App
	httpServer *http.Server
	logg       ServerLogger
}

type ServerLogger interface {
	Info(msg string, attrs ...slog.Attr)
}

// Application интерфейс для работы с событиями.
type Application interface {
	CreateEvent(ctx context.Context, e model.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (model.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

func NewServer(app *app.App, conf Config) *Server {
	return &Server{
		app: app,
		httpServer: &http.Server{
			Addr:        fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			ReadTimeout: 10 * time.Second,
			IdleTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.app.Logger.Info("Starting HTTP server...")
	s.registerLogger(s.registerRoutes())

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.app.Logger.Error("HTTP server ListenAndServe: " + err.Error())
		}
	}()

	s.app.Logger.Info("HTTP Waiting ctx done")
	<-ctx.Done()
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.app.Logger.Info("Stopping HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.app.Logger.Error("HTTP server Shutdown: " + err.Error())
		return err
	}

	s.app.Logger.Info("HTTP server stopped")
	return nil
}

func (s *Server) registerRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("POST /events", s.createEventHandler)
	mux.HandleFunc("GET /events/{id}", s.getEventHandler)
	mux.HandleFunc("DELETE /events/{id}", s.deleteEventHandler)
	return mux
}

func (s *Server) registerLogger(mux *http.ServeMux) {
	s.httpServer.Handler = loggingMiddleware(mux, s.logg)
}
