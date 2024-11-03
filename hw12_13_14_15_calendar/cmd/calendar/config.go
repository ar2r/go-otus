package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	App      AppConf      `toml:"app"`
	Logger   LoggerConf   `toml:"logger"`
	Database DatabaseConf `toml:"database"`
}

type AppConf struct {
	Env     string `toml:"env"`
	Debug   bool   `toml:"debug"`
	Storage string `toml:"storage"`
}

type LoggerConf struct {
	Level    string `toml:"level"`
	Channel  string `toml:"channel"`
	Filename string `toml:"filename"`
}

type DatabaseConf struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Schema   string `toml:"schema"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
