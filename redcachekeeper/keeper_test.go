package redcachekeeper

import (
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"

	"github.com/alicebob/miniredis"
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

func TestGetLockStore(t *testing.T) {
	// Initialize new cache keeper
	k := NewKeeper()

	m, err := miniredis.Run()
	assert.NoError(t, err)

	r := newRedisConn(m.Addr())
	k.SetConnectionPool(r)
	k.SetLockConnectionPool(r)
	k.SetWaitTime(1 * time.Second) // override wait time to 1 second

	testKey := "test-key"

	// It should return mutex when no other process is locking the process
	res, mu, err := k.GetOrLock(testKey)
	assert.Nil(t, res)
	assert.NoError(t, err)
	assert.NotNil(t, mu)

	// It should wait, and return an error while waiting for cached item ready
	res2, mu2, err2 := k.GetOrLock(testKey)
	assert.Nil(t, res2)
	assert.Nil(t, mu2)
	assert.Error(t, err2)

	// It should get response when mutex lock unlocked and cache item ready
	item := NewItem(testKey, "test-response")
	err = k.Store(mu, item)
	assert.NoError(t, err)

	res2, mu2, err2 = k.GetOrLock(testKey)
	assert.EqualValues(t, "test-response", res2)
	assert.Nil(t, mu2)
	assert.NoError(t, err2)
}

func TestPurge(t *testing.T) {
	// Initialize new cache keeper
	k := NewKeeper()

	m, err := miniredis.Run()
	assert.NoError(t, err)

	r := newRedisConn(m.Addr())
	k.SetConnectionPool(r)
	k.SetLockConnectionPool(r)

	// It should purge keys match with the matchstring while leaving the rest untouched
	testKeys := map[string]interface{}{
		"story:1234:comment:4321": nil,
		"story:1234:comment:4231": nil,
		"story:1234:comment:4121": nil,
		"story:1234:comment:4421": nil,
		"story:1234:comment:4521": nil,
		"story:1234:comment:4021": nil,
		"story:2000:comment:3021": "anything",
		"story:2000:comment:3421": "anything",
		"story:2000:comment:3231": "anything",
	}

	for key := range testKeys {
		_, mu, err := k.GetOrLock(key)
		assert.NoError(t, err)

		err = k.Store(mu, NewItem(key, "anything"))
		assert.NoError(t, err)
	}

	err = k.Purge("story:1234:*")
	assert.NoError(t, err)

	for key, value := range testKeys {
		res, _, err := k.GetOrLock(key)
		assert.NoError(t, err)
		assert.EqualValues(t, value, res)
	}
}

func TestDecideCacheTTL(t *testing.T) {
	k := &keeper{
		defaultTTL:   defaultTTL,
		lockDuration: defaultLockDuration,
		lockTries:    defaultLockTries,
		waitTime:     defaultWaitTime,
	}

	testKey := "test-key"

	// It should use keeper's default TTL when new cache item didn't specify the TTL
	i := NewItem(testKey, nil)
	assert.Equal(t, k.defaultTTL.Seconds(), k.decideCacheTTL(i))

	// It should use specified TTL when new cache item specify the TTL
	i2 := NewItemWithCustomTTL(testKey, nil, 10*time.Second)
	assert.Equal(t, i2.GetTTLFloat64(), k.decideCacheTTL(i))
}
