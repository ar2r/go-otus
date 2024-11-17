package config

type LoggerConf struct {
	Level    string `toml:"level"`
	Channel  string `toml:"channel"`
	Filename string `toml:"filename"`
}
