package myslog

import (
	"bytes"
	"log/slog"
	"testing"
)

type Logger struct {
	level  string
	logger *slog.Logger
}

func (l *Logger) Debug(msg string) {
	if l.level == DEBUG {
		l.logger.Debug(msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.level == DEBUG || l.level == INFO {
		l.logger.Info(msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level == DEBUG || l.level == INFO || l.level == WARN {
		l.logger.Warn(msg)
	}
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func newTestLogger(level string) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{})
	logger := slog.New(handler)
	return &Logger{level: level, logger: logger}, &buf
}

func TestLogger_Info(t *testing.T) {
	l, buf := newTestLogger(INFO)

	l.Info("info message")
	if !bytes.Contains(buf.Bytes(), []byte("info message")) {
		t.Errorf("expected info message to be logged")
	}

	buf.Reset()
	l.Debug("debug message")
	if buf.Len() != 0 {
		t.Errorf("expected debug message not to be logged")
	}
}

func TestLogger_Warn(t *testing.T) {
	l, buf := newTestLogger(WARN)

	l.Warn("warn message")
	if !bytes.Contains(buf.Bytes(), []byte("warn message")) {
		t.Errorf("expected warn message to be logged")
	}

	buf.Reset()
	l.Info("info message")
	if buf.Len() != 0 {
		t.Errorf("expected info message not to be logged")
	}
}

func TestLogger_Error(t *testing.T) {
	l, buf := newTestLogger(ERROR)

	l.Error("error message")
	if !bytes.Contains(buf.Bytes(), []byte("error message")) {
		t.Errorf("expected error message to be logged")
	}

	buf.Reset()
	l.Warn("warn message")
	if buf.Len() != 0 {
		t.Errorf("expected warn message not to be logged")
	}
}
