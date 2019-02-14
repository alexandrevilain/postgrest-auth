package config

import (
	"github.com/kelseyhightower/envconfig"
)

// API is the api-related configuration struct
type API struct {
	Port           int    `default:"3001"`
	ResetToken     string `default:"supersecret"`
	AllowedDomains []string
}

// Links is the links-related configuration struct
type Links struct {
	Reset   string `default:"http://localhost/reset/%v"`
	Confirm string `default:"http://localhost/confirm/%v?token=%v"`
}

// OAuth2 State is the same string that was defined to retrive the access token
type OAuth2 struct {
	State string `default:"random-state"`
}

// JWT is the jwt-related configuration struct
type JWT struct {
	Exp    int    `default:"24"`
	Secret string `default:"supersecret"`
}

// DB is the database-related configuration struct
type DB struct {
	ConnectionString string `default:"postgres://user:pass@localhost/app"`
	Roles            struct {
		Anonymous string `default:"anonymous"`
		User      string `default:"normal_user"`
	}
}

// App is the app-related configuration struct
// App is referring the the whole app where the service is deployed
type App struct {
	Name string `default:""`
	Link string `default:""`
	Logo string `default:""`
}

// Email is the email-related configuration struct
type Email struct {
	From string
	Host string
	Port int
	Auth struct {
		User string
		Pass string
	}
}

// Config represents the global config of the service
type Config struct {
	API    API
	DB     DB
	Email  Email
	JWT    JWT
	Links  Links
	App    App
	OAuth2 OAuth2
}

// LoadFromEnv loads the configuration file and populate the Config struct
func LoadFromEnv() (Config, error) {
	var config Config
	err := envconfig.Process("POSTGREST_AUTH", &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
