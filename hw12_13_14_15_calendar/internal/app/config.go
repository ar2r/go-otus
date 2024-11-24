package app

type Config struct {
	Env     string `toml:"env"`
	Debug   bool   `toml:"debug"`
	Storage string `toml:"storage"` // memory, postgres, etc.
	Queue   string `toml:"queue"`   // rabbitmq, kafka, etc.
}
