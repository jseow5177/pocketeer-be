package config

type Option struct {
	LogLevel   string
	ConfigFile string
	Port       int
}

func NewOptions() *Option {
	return &Option{
		LogLevel:   LogLevelDebug,
		ConfigFile: "",
		Port:       9090,
	}
}
