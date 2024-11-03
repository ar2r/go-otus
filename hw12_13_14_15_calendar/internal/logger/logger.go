package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	level    string
	channel  string
	filename string
	logger   *log.Logger
}

func New(level string, channel string, filename string) *Logger {
	var writer io.Writer
	var err error

	switch channel {
	case "file":
		writer, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file")
		}
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		log.Fatal("Invalid log channel")
	}

	return &Logger{
		level:    level,
		channel:  channel,
		filename: filename,
		logger:   log.New(writer, "", log.LstdFlags),
	}
}

func (l *Logger) Report() {
	l.Debug("Logger level: " + l.level + " (" + l.channel + ")")
	if l.filename != "" {
		l.Debug("Logger filename: " + l.filename)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level == "debug" {
		l.logger.Println("⚙️ DEBUG: " + msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.level == "debug" || l.level == "info" {
		l.logger.Println("🔵 INFO: " + msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level == "debug" || l.level == "info" || l.level == "warn" {
		l.logger.Println("🟡 WARN: " + msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.level == "debug" || l.level == "info" || l.level == "warn" || l.level == "error" {
		l.logger.Println("🔴 ERROR: " + msg)
	}
}
