package calendar

import "time"

type Config struct {
	Env             string        `toml:"env" default:"prod"`
	Debug           bool          `toml:"debug" default:"false"`
	Storage         string        `toml:"storage" default:"memory"` // memory, postgres, etc.
	Queue           string        `toml:"queue" default:"rabbitmq"` // rabbitmq, kafka, etc.
	CleanupDuration time.Duration `toml:"cleanup_duration" default:"100d"`
}
