package config

type Config struct {
	Server  *Server  `json:"server"`
	Mongo   *Mongo   `json:"mongo"`
	Tokens  *Tokens  `json:"token"`
	FinnHub *FinnHub `json:"finnhub"`
	IEX     *IEX     `json:"iex"`
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
	//uri := "mongodb+srv://%s:%s@%s/"
	//return fmt.Sprintf(uri, m.Username, m.Password, m.Host)
	return "mongodb://localhost:27017"
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
	Token string `json:"token"`
}

type IEX struct {
	BaseURL          string   `json:"base_url"`
	Token            string   `json:"token"`
	SupportedRegions []string `json:"supported_regions"`
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
		Tokens: &Tokens{
			AccessToken: &Token{
				ExpiresIn: 3600,
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
			Token: "cifs8bpr01qhvakk86n0cifs8bpr01qhvakk86ng",
		},
		IEX: &IEX{
			BaseURL:          "https://api.iex.cloud/v1",
			Token:            "pk_c59529d329c04cdda7708d9b7030f29f",
			SupportedRegions: []string{"US"},
		},
	}
}
