package config

type LoggerConfig struct {
	Level    string `toml:"level"`
	Channel  string `toml:"channel"`
	Filename string `toml:"filename"`
}
