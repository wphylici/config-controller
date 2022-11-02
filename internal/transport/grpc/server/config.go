package server

type Config struct {
	Network  string `toml:"network"`
	BindAddr string `toml:"bind_addr"`
}

func NewConfig() *Config {
	return &Config{
		Network:  "tcp",
		BindAddr: ":8080",
	}
}
