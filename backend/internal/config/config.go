package config

import (
	"os"
)

type Postgres struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     string
}

type HTTPServer struct {
	Address string
	Port    string
}

type Config struct {
	Postgres   Postgres
	HTTPServer HTTPServer
}

func MustLoad() *Config {
	if len(os.Args) == 2 && os.Args[1] == "local" { //todo
		return &Config{
			Postgres: Postgres{
				Port:     "5432",
				Host:     "localhost",
				DB:       "contacts_list",
				User:     "admin",
				Password: "qwerty123",
			},
			HTTPServer: HTTPServer{
				Port:    "8081",
				Address: "127.0.0.1",
			},
		}
	}

	return &Config{
		Postgres: Postgres{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
		},
		HTTPServer: HTTPServer{
			Address: os.Getenv("HTTP_SERVER_ADDRESS"),
			Port:    os.Getenv("HTTP_SERVER_PORT"),
		},
	}
}
