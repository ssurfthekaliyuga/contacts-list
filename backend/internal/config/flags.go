package config

import (
	"flag"
)

var (
	defaultConfigPath = "config.yaml"
	warnEnv           = "(please, do not use .env files in production environment it is not secure)"
	envPath           = flag.String("env", "", "Path to .env credentials file")
)

func init() {
	flag.Parse()
}
