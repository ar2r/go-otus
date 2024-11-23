package app

type Config struct {
	Env     string `toml:"env"`
	Debug   bool   `toml:"debug"`
	Storage string `toml:"storage"`
}
