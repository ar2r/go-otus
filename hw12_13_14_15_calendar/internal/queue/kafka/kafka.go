package kafka

import (
	"log/slog"
)

type Service struct {
	// some fields...
}

func NewService(logg *slog.Logger, _ Config) *Service {
	logg.Error("kafka service connector not implemented")
	return &Service{}
}
