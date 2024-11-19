package config

type GrpcServerConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
