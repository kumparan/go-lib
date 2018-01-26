package testenv

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/lab46/example/pkg/log"
)

// EnvConfig to store testing environment configuration
var EnvConfig struct {
	PostgresDSN  string `envconfig:"POSTGRES_DSN" default:"postgres://exampleapp:exampleapp@localhost:5432?sslmode=disable"`
	MySQLDSN     string `envconfig:"MYSQL_DSN"`
	RedisAddress string `envconfig:"REDIS_ADDRESS"`
}

func init() {
	err := envconfig.Process("", &EnvConfig)
	if err != nil {
		log.Errorf("Failed to load test env: %s", err.Error())
	}
}
