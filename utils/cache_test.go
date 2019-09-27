package utils

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/kumparan/cacher"
	"github.com/kumparan/tapao"
	"github.com/stretchr/testify/assert"
)

func newRedisConn(url string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     100,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", url)
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

func TestSetNotFoundCache(t *testing.T) {
	k := cacher.NewKeeper()
	m, err := miniredis.Run()

	assert.NoError(t, err)

	r := newRedisConn(m.Addr())
	k.SetConnectionPool(r)
	k.SetLockConnectionPool(r)
	k.SetWaitTime(1 * time.Second)

	testKey := "test-key"

	// set a not founc cache
	err = SetNotFoundCache(k, testKey, 5*time.Minute)
	assert.NoError(t, err)
	type testStruct struct{}
	var data *testStruct
	reply, mu, err := k.GetOrLock(testKey)

	_ = tapao.Unmarshal(reply.([]byte), data)

	assert.Nil(t, data)
	assert.Nil(t, mu)
	assert.Nil(t, err)
}
