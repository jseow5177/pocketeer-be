package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type Config struct {
	RateLimits          map[string]*RateLimit `json:"rate_limits"`
	ServerAdmin         *ServerAdmin          `json:"server_admin"`
	Mongo               *Mongo                `json:"mongo"`
	Tokens              *Tokens               `json:"tokens"`
	FinnHub             *FinnHub              `json:"finnhub"`
	ExchangeRateHost    *ExchangeRateHost     `json:"exchange_rate_host"`
	QuoteMemCache       *MemCache             `json:"quote_mem_cache"`
	FeedbackGoogleSheet *GoogleSheet          `json:"feedback_google_sheet"`
	OTPMemCache         *MemCache             `json:"otp_mem_cache"`
	Brevo               *Brevo                `json:"brevo"`
	Gmail               *Gmail                `json:"gmail"`
	Global              *Global               `json:"global"`
}

type Global struct {
	UseGmail bool `json:"use_gmail"`
}

type ServerAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RateLimit struct {
	RPS   float64 `json:"rps"`
	Burst int64   `json:"burst"`
}

type MemCache struct {
	ExpiryTime      string `json:"expiry_time"`
	CleanUpInterval string `json:"clean_up_interval"`
}

type Mongo struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Database        string `json:"database"`
	ConnectTimeout  string `json:"connect_timeout"`
	MaxConnIdleTime string `json:"max_conn_idle_time"`
	MinConnPoolSize int    `json:"min_conn_pool_size"`
	MaxConnPoolSize int    `json:"max_conn_pool_size"`
	Timeout         string `json:"timeout"`
	RetryReads      *bool  `json:"retry_reads"`
	RetryWrites     *bool  `json:"retry_writes"`
}

type GoogleSheet struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	Type        string `json:"type"`
	SheetID     string `json:"sheet_id"`
	WriteRange  string `json:"write_range"`
}

func (m *Mongo) String() string {
	dsn := &url.URL{
		Scheme:     "mongodb+srv",
		User:       url.UserPassword(m.Username, m.Password),
		Host:       m.Host,
		Path:       fmt.Sprintf("/%s", m.Database),
		ForceQuery: false,
	}

	q := dsn.Query()

	if m.ConnectTimeout != "" {
		if t, err := time.ParseDuration(m.ConnectTimeout); err == nil {
			q.Set("connectTimeoutMS", fmt.Sprint(t.Milliseconds()))
		}
	}

	if m.Timeout != "" {
		if t, err := time.ParseDuration(m.Timeout); err == nil {
			q.Set("timeoutms", fmt.Sprint(t.Milliseconds()))
		}
	}

	if m.MaxConnIdleTime != "" {
		if t, err := time.ParseDuration(m.MaxConnIdleTime); err == nil {
			q.Set("maxIdleTimeMS", fmt.Sprint(t.Milliseconds()))
		}
	}

	if m.RetryReads != nil {
		q.Set("retryReads", fmt.Sprint(*m.RetryReads))
	}

	if m.RetryWrites != nil {
		q.Set("retryWrites", fmt.Sprint(*m.RetryWrites))
	}

	q.Set("minPoolSize", fmt.Sprint(m.MinConnPoolSize))

	q.Set("maxPoolSize", fmt.Sprint(m.MaxConnPoolSize))

	q.Set("authSource", "admin")

	dsn.RawQuery = q.Encode()

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
}

type FinnHub struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

type ExchangeRateHost struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

type Brevo struct {
	APIKey string `json:"api_key"`
}

type Gmail struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewConfig() *Config {
	return &Config{
		RateLimits: make(map[string]*RateLimit),
		ServerAdmin: &ServerAdmin{
			Username: "admin",
			Password: "",
		},
		Mongo: &Mongo{
			Username:        "pocketeer-test",
			Password:        "",
			Host:            "pocketeer-test.aepotln.mongodb.net",
			Database:        "pocketeer",
			ConnectTimeout:  "5s",
			MaxConnIdleTime: "60m",
			Timeout:         "5s",
			MinConnPoolSize: 20,
			MaxConnPoolSize: 100,
			RetryReads:      goutil.Bool(true),
			RetryWrites:     goutil.Bool(true),
		},
		Tokens: &Tokens{
			AccessToken: &Token{
				ExpiresIn: 31_536_000, // 365 days
				Issuer:    "pocketeer_be",
				Secret:    "",
			},
			RefreshToken: &Token{
				ExpiresIn: 3600,
				Issuer:    "pocketeer_be",
				Secret:    "",
			},
		},
		FinnHub: &FinnHub{
			BaseURL: "https://finnhub.io/api/v1",
			Token:   "",
		},
		QuoteMemCache: &MemCache{
			ExpiryTime:      "15m",
			CleanUpInterval: "20m",
		},
		FeedbackGoogleSheet: &GoogleSheet{
			ClientEmail: "",
			PrivateKey:  "",
			Type:        "",
			SheetID:     "",
			WriteRange:  "",
		},
		OTPMemCache: &MemCache{
			ExpiryTime:      "10m",
			CleanUpInterval: "15m",
		},
		Brevo: &Brevo{
			APIKey: "",
		},
		Gmail: &Gmail{
			Name:     "Bytewise",
			Host:     "smtp.gmail.com",
			Port:     587,
			Email:    "",
			Password: "",
		},
		ExchangeRateHost: &ExchangeRateHost{
			BaseURL: "https://api.currencyapi.com/v3",
			APIKey:  "",
		},
		Global: &Global{
			UseGmail: true,
		},
	}
}

func (c *Config) Subscribe(ctx context.Context, configFile string) error {
	if configFile == "" {
		log.Ctx(ctx).Warn().Msgf("empty config file")
		return nil
	}

	f, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Ctx(ctx).Warn().Msgf("config file does not exist, file path: %v", configFile)
			return nil
		}
		return err
	}
	defer f.Close()

	p := json.NewDecoder(f)
	if err := p.Decode(&c); err != nil {
		return err
	}

	return nil
}
