package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// NetworkTCP for tcp
const NetworkTCP = "tcp"

type (
	// Config of redis
	Config struct {
		Address   string `yaml:"address"`
		MaxIdle   int    `yaml:"maxidle"`
		MaxActive int    `yaml:"maxactive"`
		Timeout   int    `yaml:"timeout"`
	}
	// Redis struct
	Redis struct {
		Pool *redigo.Pool
	}
)

// Init redis connection
func Init(cfg Config) *Redis {
	pool := &redigo.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: time.Duration(cfg.Timeout) * time.Second,
		Wait:        true,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial(NetworkTCP, cfg.Address)
		},
	}
	r := &Redis{Pool: pool}
	return r
}

func (r *Redis) Err() error {
	return r.Pool.Get().Err()
}
