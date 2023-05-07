package config

import "fmt"

type Config struct {
	Server *Server `json:"server"`
	Mongo  *Mongo  `json:"mongo"`
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

type Mongo struct {
	Username string `json:"username"`
	Password string `json:"passowrd"`
	Host     string `json:"host"`
	DBName   string `json:"db_name"`
}

func (m *Mongo) String() string {
	uri := "mongodb+srv://%s:%s@%s/"
	return fmt.Sprintf(uri, m.Username, m.Password, m.Host)
}

func NewConfig() *Config {
	return &Config{
		Server: &Server{
			LogLevel:   "DEBUG",
			Port:       9090,
			RateLimits: make(map[string]*RateLimit),
		},
		Mongo: &Mongo{
			Username: "pocketeer-test",
			Password: "eTSvssKfSWCzRylk",
			Host:     "mongodb-test.djhnkbj.mongodb.net",
			DBName:   "pocketeer",
		},
	}
}
