package config

type RestServerConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
