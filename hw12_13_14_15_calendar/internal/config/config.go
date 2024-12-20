package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/database"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue/kafka"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
)

type Config struct {
	App        calendar.Config   `toml:"app"`
	HTTPServer httpserver.Config `toml:"http"`
	GRPCServer grpcserver.Config `toml:"grpc"`
	Database   database.Config   `toml:"database"`
	Logger     LoggerConfig      `toml:"logger"`
	RabbitMQ   rabbit.Config     `toml:"rabbitmq"`
	Kafka      kafka.Config      `toml:"kafka"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
