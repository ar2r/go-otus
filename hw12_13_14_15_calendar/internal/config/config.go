package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	App        AppConf        `toml:"app"`
	RestServer RestServerConf `toml:"rest"`
	GrpcServer GrpcServerConf `toml:"grpc"`
	Logger     LoggerConf     `toml:"logger"`
	Database   DatabaseConf   `toml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
