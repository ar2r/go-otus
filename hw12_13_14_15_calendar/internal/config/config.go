package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/database"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/http"
)

type Config struct {
	App        app.Config        `toml:"app"`
	HttpServer httpserver.Config `toml:"http"`
	GrpcServer grpcserver.Config `toml:"grpc"`
	Database   database.Config   `toml:"database"`
	Logger     LoggerConf        `toml:"logger"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
