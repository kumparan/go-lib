package redis

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var (
	// Pool -> redis connection pool
	Pool *redigo.Pool
)

// Init initialize Redis Pool
func Init(redisHost string) {

	Pool = newPool(redisHost)
	cleanupHook()
}

func newPool(server string) *redigo.Pool {

	return &redigo.Pool{

		MaxIdle:     100,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
