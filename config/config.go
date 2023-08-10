package config

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	Server        *Server   `json:"server"`
	Mongo         *Mongo    `json:"mongo"`
	Tokens        *Tokens   `json:"token"`
	FinnHub       *FinnHub  `json:"finnhub"`
	QuoteMemCache *MemCache `json:"quote_mem_cache"`
	Brevo         *Brevo    `json:"brevo"`
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

type MemCache struct {
	ExpiryTime      string `json:"expiry_time"`
	CleanUpInterval string `json:"clean_up_interval"`
}

type Mongo struct {
	Username       string `json:"username"`
	Password       string `json:"passowrd"`
	Host           string `json:"host"`
	Database       string `json:"database"`
	ConnectTimeout string `json:"connect_timeout"`
}

func (m *Mongo) String() string {
	dsn := &url.URL{
		Scheme: "mongodb+srv",
		User:   url.UserPassword(m.Username, m.Password),
		Host:   m.Host,
	}

	q := dsn.Query()

	if m.ConnectTimeout != "" {
		if t, err := time.ParseDuration(m.ConnectTimeout); err == nil {
			q.Set("connectTimeoutMS", fmt.Sprint(t.Milliseconds()))
		}
	}

	return dsn.String()
}

type Token struct {
	Secret    string `json:"secret"`
	ExpiresIn int64  `json:"expires_in"` // second
	Issuer    string `json:"issuer"`
}

type Tokens struct {
	AccessToken  *Token `json:"access_token"`
	RefreshToken *Token `json:"refresh_token"`
	EmailToken   *Token `json:"email_token"`
}

type FinnHub struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

type Brevo struct {
	APIKey string `json:"api_key"`
}

func NewConfig() *Config {
	return &Config{
		Server: &Server{
			LogLevel:   "DEBUG",
			Port:       9090,
			RateLimits: make(map[string]*RateLimit),
		},
		Mongo: &Mongo{
			Username:       "pocketeer-test",
			Password:       "eTSvssKfSWCzRylk",
			Host:           "mongodb-test.djhnkbj.mongodb.net",
			Database:       "pocketeer",
			ConnectTimeout: "5s",
		},
		Tokens: &Tokens{
			AccessToken: &Token{
				ExpiresIn: 31_536_000, // 365 days
				Issuer:    "pocketeer_be",
				Secret:    "%5jJclw22Sa91k9V4N11H^zGXkc0jw",
			},
			RefreshToken: &Token{
				ExpiresIn: 3600,
				Issuer:    "pocketeer_be",
				Secret:    "@w8DlsuWfSlg25W0#qbZ5CpGq#MNlB",
			},
			EmailToken: &Token{
				ExpiresIn: 604_800,
				Issuer:    "pocketeer_be",
				Secret:    "hWjtP3Hzizi2VpaZSpchriMSJ7EqocDM",
			},
		},
		FinnHub: &FinnHub{
			BaseURL: "https://finnhub.io/api/v1",
			Token:   "cifs8bpr01qhvakk86n0cifs8bpr01qhvakk86ng",
		},
		QuoteMemCache: &MemCache{
			ExpiryTime:      "15m",
			CleanUpInterval: "20m",
		},
		Brevo: &Brevo{
			APIKey: "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",
		},
	}
}
