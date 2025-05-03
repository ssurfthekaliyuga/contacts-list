package config

import (
	"time"
)

type Config struct {
	LoggerConf     Logger     `mapstructure:"logger"`
	PostgresConf   Postgres   `mapstructure:"postgres"`
	HTTPServerConf HTTPServer `mapstructure:"http_server"`
}

type Logger struct {
	Level       int    `mapstructure:"level"`
	AddSource   bool   `mapstructure:"add_source"`
	HandlerType string `mapstructure:"handler_type"`
}

type Postgres struct {
	ConnString string `mapstructure:"conn_string"`
}

type HTTPServer struct {
	Host            string                `mapstructure:"host"`
	Port            string                `mapstructure:"port"`
	ShutdownTimeout time.Duration         `mapstructure:"shutdown_timeout"`
	EndpointsV1     HTTPServerEndpointsV1 `mapstructure:"endpoints_v1"`
}

type HTTPServerEndpointsV1 struct {
	Cors HTTPServerCors `mapstructure:"cors"`
}

type HTTPServerCors struct {
	AllowOrigins string `mapstructure:"allow_origins"`
	AllowMethods string `mapstructure:"allow_methods"`
	AllowHeaders string `mapstructure:"allow_headers"`
}
