package config

type ServerConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
