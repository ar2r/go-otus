package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &Logger{level: "debug", logger: logger}

	l.Debug("debug message")
	if !bytes.Contains(buf.Bytes(), []byte("DEBUG: debug message")) {
		t.Errorf("expected debug message to be logged")
	}
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &Logger{level: "info", logger: logger}

	l.Info("info message")
	if !bytes.Contains(buf.Bytes(), []byte("INFO: info message")) {
		t.Errorf("expected info message to be logged")
	}

	buf.Reset()
	l.Debug("debug message")
	if buf.Len() != 0 {
		t.Errorf("expected debug message not to be logged")
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &Logger{level: "warn", logger: logger}

	l.Warn("warn message")
	if !bytes.Contains(buf.Bytes(), []byte("WARN: warn message")) {
		t.Errorf("expected warn message to be logged")
	}

	buf.Reset()
	l.Info("info message")
	if buf.Len() != 0 {
		t.Errorf("expected info message not to be logged")
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &Logger{level: "error", logger: logger}

	l.Error("error message")
	if !bytes.Contains(buf.Bytes(), []byte("ERROR: error message")) {
		t.Errorf("expected error message to be logged")
	}

	buf.Reset()
	l.Warn("warn message")
	if buf.Len() != 0 {
		t.Errorf("expected warn message not to be logged")
	}
}
