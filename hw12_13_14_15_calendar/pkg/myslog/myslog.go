package myslog

import (
	"log"
	"log/slog"
	"os"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

func New(level string, channel string, filename string) *slog.Logger {
	var slogHandler slog.Handler
	var slogLevel slog.Level

	switch level {
	case DEBUG:
		slogLevel = slog.LevelDebug
	case INFO:
		slogLevel = slog.LevelInfo
	case WARN:
		slogLevel = slog.LevelWarn
	case ERROR:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelError
	}

	options := &slog.HandlerOptions{
		Level: slogLevel,
	}

	switch channel {
	case "file":
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			log.Fatal("Failed to open log file")
		}
		slogHandler = slog.NewTextHandler(file, options)
	case "stdout":
		slogHandler = slog.NewTextHandler(os.Stdout, options)
	case "stderr":
		slogHandler = slog.NewTextHandler(os.Stderr, options)
	default:
		log.Fatal("Invalid log channel")
	}

	return slog.New(slogHandler)
}
