package config

type Option struct {
	LogLevel   string
	ConfigFile string
	Port       int
}

func NewOptions() *Option {
	return &Option{
		LogLevel:   LogLevelDebug,
		ConfigFile: "./bin/config-test.json",
		Port:       9090,
	}
}
