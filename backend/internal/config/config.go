package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Postgres   Postgres   `mapstructure:"postgres"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Logger     Logger     `mapstructure:"logger"`
}

type Postgres struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type HTTPServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Logger struct {
	Level       int    `mapstructure:"level"`
	AddSource   bool   `mapstructure:"add_source"`
	HandlerType string `mapstructure:"handler_type"`
}

func Read() (*Config, error) {
	if err := godotenv.Load(); err != nil { //todo delete it when if i will start using dev container
		return nil, fmt.Errorf("cannot read .env file (you must not use .env file in production): %w", err)
	}

	vp := viper.New()

	vp.SetConfigFile(configPath())
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var conf Config
	if err := vp.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config from file: %w", err)
	}

	return conf.expandEnv()
}

func (c *Config) expandEnv() (*Config, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal config to json in order to expand env varibles: %w", err)
	}

	expanded := os.ExpandEnv(string(bytes))
	_ = json.Unmarshal([]byte(expanded), c)

	return c, nil
}

func configPath() string {
	var (
		pathDefault = "config.yaml"
		pathFlag    = flag.String("config", "", "path to yaml config file")
		pathEnv     = os.Getenv("CONFIG")
	)

	flag.Parse()

	switch {
	case *pathFlag != "":
		return *pathFlag
	case pathEnv != "":
		return pathEnv
	default:
		return pathDefault
	}
}
