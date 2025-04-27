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
	Postgres   string     `mapstructure:"postgres"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Logger     Logger     `mapstructure:"logger"`
}

type HTTPServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Logger struct { //todo time format
	Level       int    `mapstructure:"level"`
	AddSource   bool   `mapstructure:"add_source"`
	HandlerType string `mapstructure:"handler_type"`
}

var (
	defaultConfigPath = "config.yaml"
	warnEnv           = "(please, do not use .env files in production environment it is not secure)"
	envPath           = flag.String("env", "", "Path to .env credentials file")
)

func init() {
	flag.Parse()
}

func Read() (*Config, error) {
	if envPath != nil && *envPath != "" {
		err := godotenv.Load(*envPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read .env file (%s) (%s): %w", *envPath, warnEnv, err)
		}
	}

	var configPath string
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	} else {
		configPath = defaultConfigPath
	}

	vp := viper.New()

	vp.SetConfigFile(configPath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config file (%s): %w", configPath, err)
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

	if err := json.Unmarshal([]byte(expanded), c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config from json in order to expand env varibles: %w", err)
	}

	return c, nil
}
