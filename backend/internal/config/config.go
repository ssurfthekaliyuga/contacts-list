package config

import (
	"contacts-list/internal/common/logger"
	"contacts-list/internal/primary/rest/endpoints"
	"contacts-list/internal/primary/rest/fiber"
	"contacts-list/internal/secondary/postgres"
	"contacts-list/pkg/sl"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

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

func (c *Config) Logger() logger.Config {
	return logger.Config{
		AddSource:   c.LoggerConf.AddSource,
		Level:       sl.Level(c.LoggerConf.Level),
		HandlerType: logger.HandlerType(c.LoggerConf.HandlerType),
	}
}

func (c *Config) Postgres() postgres.Config {
	return postgres.Config{
		ConnString: c.PostgresConf.ConnString,
	}
}

func (c *Config) FiberServer() fiber.Config {
	return fiber.Config{
		Host: c.HTTPServerConf.Host,
		Port: c.HTTPServerConf.Port,
	}
}

func (c *Config) EndpointsV1() endpoints.V1Config {
	return endpoints.V1Config{
		Cors: cors.Config{
			Next:             nil,
			AllowOriginsFunc: nil,
			AllowOrigins:     c.HTTPServerConf.EndpointsV1.Cors.AllowOrigins,
			AllowMethods:     c.HTTPServerConf.EndpointsV1.Cors.AllowMethods,
			AllowHeaders:     c.HTTPServerConf.EndpointsV1.Cors.AllowHeaders,
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           0,
		},
	}
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
