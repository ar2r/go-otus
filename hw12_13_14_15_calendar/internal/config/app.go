package config

type AppConf struct {
	Env     string `toml:"env"`
	Debug   bool   `toml:"debug"`
	Storage string `toml:"adapters"`
}
