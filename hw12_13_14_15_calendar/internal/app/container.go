package app

import "log/slog"

var (
	logger *slog.Logger
)

func SetLogger(l *slog.Logger) {
	logger = l
}

func GetLogger() *slog.Logger {
	return logger
}
