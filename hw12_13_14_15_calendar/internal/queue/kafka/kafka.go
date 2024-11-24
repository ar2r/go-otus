package kafka

import (
	"log/slog"
)

type Service struct {
	// some fields...
}

func NewService(logg *slog.Logger, conf Config) *Service {
	logg.Error("kafka service connector not implemented")
	return &Service{}
}
