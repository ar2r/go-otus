package grpcserver

type Config struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
