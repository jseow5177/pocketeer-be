package config

type Config struct {
	Server    *Server    `json:"server"`
	Firestore *FireStore `json:"firestore"`
}

type RateLimit struct {
	RPS   float64 `json:"rps"`
	Burst int64   `json:"burst"`
}

type Server struct {
	LogLevel   string                `json:"log_level"`
	Port       int                   `json:"port"`
	RateLimits map[string]*RateLimit `json:"rate_limits"`
}

type FireStore struct{}

func NewConfig() *Config {
	return &Config{
		Server: &Server{
			LogLevel:   "DEBUG",
			Port:       9090,
			RateLimits: make(map[string]*RateLimit),
		},
		Firestore: &FireStore{},
	}
}
