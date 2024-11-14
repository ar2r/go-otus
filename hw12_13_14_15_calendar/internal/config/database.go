package config

type DatabaseConf struct {
	Host              string `toml:"host"`
	Port              int    `toml:"port"`
	Database          string `toml:"database"`
	Username          string `toml:"username"`
	Password          string `toml:"password"`
	Schema            string `toml:"schema"`
	SSLMode           string `toml:"ssl_mode"`
	Timezone          string `toml:"timezone"`
	TargetSessionAttr string `toml:"target_session_attr"`
}
