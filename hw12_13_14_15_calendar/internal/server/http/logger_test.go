package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockLogger struct {
	logs []string
}

func (m *MockLogger) Log(_ slog.Level, msg string, _ ...interface{}) {
	m.logs = append(m.logs, msg)
}

func (m *MockLogger) Debug(msg string, attrs ...interface{}) {
	m.Log(slog.LevelDebug, msg, attrs...)
}

func (m *MockLogger) Info(msg string, attrs ...interface{}) {
	m.Log(slog.LevelInfo, msg, attrs...)
}

func (m *MockLogger) Warn(msg string, attrs ...interface{}) {
	m.Log(slog.LevelWarn, msg, attrs...)
}

func (m *MockLogger) Error(msg string, attrs ...interface{}) {
	m.Log(slog.LevelError, msg, attrs...)
}

func TestLoggingMiddleware(t *testing.T) {
	logg := &MockLogger{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	wrappedHandler := loggingMiddleware(handler, logg)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/hello?q=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.RemoteAddr = "11.22.33.44:12345"

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if len(logg.logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logg.logs))
	}

	logEntry := logg.logs[0]
	if !strings.Contains(logEntry, "11.22.33.44") ||
		!strings.Contains(logEntry, "GET /hello?q=1 HTTP/1.1") ||
		!strings.Contains(logEntry, "200") ||
		!strings.Contains(logEntry, "Mozilla/5.0") {
		t.Errorf("log entry does not contain expected values: %s", logEntry)
	}
}
