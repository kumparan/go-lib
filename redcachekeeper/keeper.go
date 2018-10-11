package redcachekeeper

import (
	"errors"
	"time"

	"github.com/go-redsync/redsync"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jpillora/backoff"
)

const (
	// Override these when constructing the cache keeper
	defaultTTL          = 10 * time.Second
	defaultLockDuration = 1 * time.Minute
	defaultLockTries    = 1
	defaultWaitTime     = 15 * time.Second
)

type (
	// Keeper responsible for managing cache
	Keeper interface {
		GetOrLock(string) (interface{}, *redsync.Mutex, error)
		Store(*redsync.Mutex, Item) error
		Purge(string) error

		SetDefaultTTL(time.Duration)
		SetConnectionPool(*redigo.Pool)
		SetLockConnectionPool(*redigo.Pool)
		SetLockDuration(time.Duration)
		SetLockTries(int)
		SetWaitTime(time.Duration)
		SetDisableCaching(bool)
	}

	keeper struct {
		connPool       *redigo.Pool
		defaultTTL     time.Duration
		waitTime       time.Duration
		disableCaching bool

		lockConnPool *redigo.Pool
		lockDuration time.Duration
		lockTries    int
	}
)

// NewKeeper :nodoc:
func NewKeeper() Keeper {
	return &keeper{
		defaultTTL:     defaultTTL,
		lockDuration:   defaultLockDuration,
		lockTries:      defaultLockTries,
		waitTime:       defaultWaitTime,
		disableCaching: false,
	}
}

// GetOrLock :nodoc:
func (k *keeper) GetOrLock(key string) (cachedItem interface{}, mutex *redsync.Mutex, err error) {
	if k.disableCaching {
		return
	}

	cachedItem, err = k.getCachedItem(key)
	if err != nil && err != redigo.ErrNil || cachedItem != nil {
		return
	}

	mutex, err = k.acquireLock(key)
	if err == nil {
		return
	}

	start := time.Now()
	for {
		b := &backoff.Backoff{
			Min:    20 * time.Millisecond,
			Max:    200 * time.Millisecond,
			Jitter: true,
		}
		cachedItem, err = k.getCachedItem(key)
		if err != nil && err != redigo.ErrNil || cachedItem != nil {
			return
		}

		elapsed := time.Since(start)
		if elapsed >= k.waitTime {
			break
		}

		time.Sleep(b.Duration())
	}

	return nil, nil, errors.New("Wait Too Long")
}

// Store :nodoc:
func (k *keeper) Store(mutex *redsync.Mutex, c Item) error {
	if k.disableCaching {
		return nil
	}
	defer mutex.Unlock()

	client := k.connPool.Get()
	defer client.Close()

	_, err := client.Do("SETEX", c.GetKey(), k.decideCacheTTL(c), c.GetValue())
	return err
}

// Purge :nodoc:
func (k *keeper) Purge(matchString string) error {
	client := k.connPool.Get()
	defer client.Close()

	keys, err := redigo.Values(client.Do("KEYS", matchString))
	if err != nil {
		return err
	}

	if keys == nil {
		return errors.New("redcachekeeper: No matching keys")
	}

	client.Send("MULTI")
	for _, k := range keys {
		client.Send("DEL", k)
	}
	_, err = client.Do("EXEC")

	return err
}

// SetDefaultTTL :nodoc:
func (k *keeper) SetDefaultTTL(d time.Duration) {
	k.defaultTTL = d
}

// SetConnectionPool :nodoc:
func (k *keeper) SetConnectionPool(c *redigo.Pool) {
	k.connPool = c
}

// SetLockConnectionPool :nodoc:
func (k *keeper) SetLockConnectionPool(c *redigo.Pool) {
	k.lockConnPool = c
}

// SetLockDuration :nodoc:
func (k *keeper) SetLockDuration(d time.Duration) {
	k.lockDuration = d
}

// SetLockTries :nodoc:
func (k *keeper) SetLockTries(t int) {
	k.lockTries = t
}

// SetWaitTime :nodoc:
func (k *keeper) SetWaitTime(d time.Duration) {
	k.waitTime = d
}

// SetDisableCaching :nodoc:
func (k *keeper) SetDisableCaching(b bool) {
	k.disableCaching = b
}

func (k *keeper) decideCacheTTL(c Item) float64 {
	if c.GetTTLFloat64() > 0 {
		return c.GetTTLFloat64()
	}

	return k.defaultTTL.Seconds()
}

func (k *keeper) acquireLock(key string) (*redsync.Mutex, error) {
	r := redsync.New([]redsync.Pool{k.lockConnPool})
	m := r.NewMutex("lock:"+key,
		redsync.SetExpiry(k.lockDuration),
		redsync.SetTries(k.lockTries))

	return m, m.Lock()
}

func (k *keeper) getCachedItem(key string) (value interface{}, err error) {
	client := k.connPool.Get()
	defer client.Close()

	return client.Do("GET", key)
}
