package config

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	Server              *Server      `json:"server"`
	Mongo               *Mongo       `json:"mongo"`
	Tokens              *Tokens      `json:"token"`
	FinnHub             *FinnHub     `json:"finnhub"`
	QuoteMemCache       *MemCache    `json:"quote_mem_cache"`
	FeedbackGoogleSheet *GoogleSheet `json:"feedback_google_sheet"`
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

type GoogleSheet struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	Type        string `json:"type"`
	SheetID     string `json:"sheet_id"`
	WriteRange  string `json:"write_range"`
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
}

type FinnHub struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
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
		},
		FinnHub: &FinnHub{
			BaseURL: "https://finnhub.io/api/v1",
			Token:   "cifs8bpr01qhvakk86n0cifs8bpr01qhvakk86ng",
		},
		QuoteMemCache: &MemCache{
			ExpiryTime:      "15m",
			CleanUpInterval: "20m",
		},
		FeedbackGoogleSheet: &GoogleSheet{
			ClientEmail: "pocketeer@dev-country-396508.iam.gserviceaccount.com",
			PrivateKey:  "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCp6icfVAR2uJgV\nuqFoesHTY7M2nsbZb7glFdsPnC+7r3vJANxVs/kQUMGZyxJAu8YUKTEGFZUjsx8K\nM5gozspzs+vfjGi9cTsNkQLqoQcICZC10RVyl4vSDR4H/ihtIda2/8XI5wYYdvQM\nb84YoepaXLqBuwo2qUsecqEvolLD81S+8S873gK4vDHF1ceUMMZnt90U2JJ/DS+Y\nNrQTNu3o/POPZv3cwmVVu6ecgXHQco26jf1n6IaEcQhRiH7lBfCNkR6W7R5CUQBV\naDY+0zlC6iZyFJv5CTmULx/FCi9KrAmJE0IHb9dR8cxtUQjgWr3HOdxnyv+2EUf5\nRYyBx1jzAgMBAAECggEACPbO+JKkp6TGdUMC1/G9/wp//m4y/G6t7trv0yuAwyyX\nYbvXY/gaODeidxwlCuA9K3we13zVAOymwxGEwgZi3ObNl14fensaS/VuC3XSWqSy\niz0MHh8Lo9nIT1CjzloeK0pkI1y40BqadnuBioAkTn8c4dI90DNvQJx4j3xHVIoQ\nUelZCLz2qxzu3YQ6AGeE//Q30A2kTa3y1h5myFazDMdlmDYmBHUp17b/lUudPB9v\n4OXPu8j23AoyyT/4D5zJX05FvomaUpowboatc0MoDmgVwAu2Ug7MFAwOrey0tGvv\nmcZIvl6EQ6KpSFuoTwXqN2d7sjzq6H4Vi9TKpi6S/QKBgQDblQgMQLeqOcVumvTc\n9HcmF7RvbI6E7Yn4hMZTSnTSYcWM9ci/DfuiyuFLVYECacKWVgsrEjSSuUKwX01z\nbz8BXsB4bNVPii8qviWqEPubqk2+29bpQ9pZg8IBOe/jW9g0gP36ooxHt9Py2xG2\nVOSjvTdUkjeb+2kTcnWNt5qN7QKBgQDGGFmvNcqKZQE0Mj0r4dfOLvl8qvVJHIH0\ni1UQW1zTj5/FsuwoVABhHLo0ReBfbXMRGO/ppp26ogvj0LUw4dpCqId978vZz3RZ\n/7mb1+ifJXcc5+l8OttrfyVbHGskDUkJ9LrMXLAhAlYNxPhzs8K3ppgDir94gSn/\ncd7hEFmmXwKBgDSjypUdYqpVRSZZ0X+yv2mLXz8i+BuX0m6Ybe8Dt3PD6pb7SQ+8\nK2oAVvg3XEEW14YwxGaj66SM8xbTEf8tWR5b96om9RAnYV1Ozjqx7Y+IyTCLBT9Q\ne+TfuD+RAxgvKWqUzc4q75Q11oKuz9U1DsbOEpicoOYs5Ci8VMAPluaVAoGBAKdv\nuY5w5wtnKEdYF4BZ2jC6X8JSNhVf9TAf/PxgXOutQBy4iQflSJTM1U6NqYK/Xj0b\nWN2jKTqw8V/T7vKsU/F3xV5bK3Ck+vF/RwFE06iA4FccZqvMi94mkAqc0KqeWBgE\nNUe7KFweP2JQFLinPnRAacjEo+ZVNoxlUT/ms//9AoGBANbO4b59KoIS6txo/Tfd\n3+/90OsD8SlqI8XlFwjwV2ypOHP4vtHn24ZvuxWGxPj/jzsp8ZP9DQBCP6OiDaXG\nffQLJ6lCdkrKIuNXLrto4Sz5aNkgFnATKNH684ukuShgUz9s1wMkt599igzvA3hW\nCH8QzmeALOfcb2Ajl0N8LeRf",
			Type:        "service_account",
			SheetID:     "1xeY7_DRK0EYxEnO_3e9Wh7kNlgEZLIq0TL5EMth6juU",
			WriteRange:  "Sheet1!A1:C2",
		},
	}
}
