package redcachekeeperv2

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
	defaultStatPrefix   = "stat:"
	defaultEnableStat   = false
)

type (
	CacheGeneratorFn func() (interface{}, error)

	// Keeper responsible for managing cache
	Keeper interface {
		GetOrLock(string) (interface{}, *redsync.Mutex, error)
		GetOrSet(string, CacheGeneratorFn, time.Duration) (interface{}, error)
		Store(*redsync.Mutex, Item) error
		Purge(string) error
		IncreaseCachedValueByOne(key string) error

		AcquireLock(string) (*redsync.Mutex, error)
		SetDefaultTTL(time.Duration)
		SetConnectionPool(*redigo.Pool)
		SetLockConnectionPool(*redigo.Pool)
		SetLockDuration(time.Duration)
		SetLockTries(int)
		SetWaitTime(time.Duration)
		SetDisableCaching(bool)
		SetStatPrefix(string)
		SetEnableStat(bool)
		ClearStats() error
	}

	keeper struct {
		connPool       *redigo.Pool
		defaultTTL     time.Duration
		waitTime       time.Duration
		disableCaching bool

		lockConnPool *redigo.Pool
		lockDuration time.Duration
		lockTries    int

		statPrefix string
		enableStat bool
	}

	statItem struct {
		Name  string
		Count int64
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
		statPrefix:     defaultStatPrefix,
		enableStat:     defaultEnableStat,
	}
}

// GetOrLock :nodoc:
func (k *keeper) GetOrLock(key string) (cachedItem interface{}, mutex *redsync.Mutex, err error) {
	if k.disableCaching {
		return
	}

	cachedItem, err = k.getCachedItem(key, true)
	if err != nil && err != redigo.ErrNil || cachedItem != nil {
		return
	}

	mutex, err = k.AcquireLock(key)
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

		if !k.isLocked(key) {
			cachedItem, err = k.getCachedItem(key, false)
			if err != nil && err != redigo.ErrNil || cachedItem != nil {
				return
			}
			return nil, nil, nil
		}

		elapsed := time.Since(start)
		if elapsed >= k.waitTime {
			break
		}

		time.Sleep(b.Duration())
	}

	return nil, nil, errors.New("wait too long")
}

// GetOrSet :nodoc:
func (k *keeper) GetOrSet(key string, fn CacheGeneratorFn, ttl time.Duration) (cachedItem interface{}, err error) {
	cachedItem, mu, err := k.GetOrLock(key)
	if err != nil {
		return
	}
	if cachedItem != nil {
		return
	}

	defer func() {
		if mu != nil {
			mu.Unlock()
		}
	}()

	cachedItem, err = fn()

	if err != nil {
		return
	}

	err = k.Store(mu, NewItemWithCustomTTL(key, cachedItem, ttl))

	return
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

// rgrge :nodoc:
func (k *keeper) Purge(matchString string) error {
	if k.disableCaching {
		return nil
	}

	client := k.connPool.Get()
	defer client.Close()

	keys, err := redigo.Values(client.Do("KEYS", matchString))
	if err != nil {
		return err
	}

	if keys == nil {
		return errors.New("redcachekeeper: No matching keys")
	}

	_, err = client.Do("DEL", keys...)

	return err
}

// ClearStats :nodoc:
func (k *keeper) ClearStats() error {
	return k.Purge(k.statPrefix + "*")
}

// IncreaseCachedValueByOne will increments the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation
func (k *keeper) IncreaseCachedValueByOne(key string) error {
	if k.disableCaching {
		return nil
	}

	client := k.connPool.Get()
	defer client.Close()

	_, err := client.Do("INCR", key)
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

// SetStatPrefix :nodoc:
func (k *keeper) SetStatPrefix(s string) {
	k.statPrefix = s
}

// SetEnableStat :nodoc:
func (k *keeper) SetEnableStat(b bool) {
	k.enableStat = b
}

// AcquireLock :nodoc:
func (k *keeper) AcquireLock(key string) (*redsync.Mutex, error) {
	r := redsync.New([]redsync.Pool{k.lockConnPool})
	m := r.NewMutex("lock:"+key,
		redsync.SetExpiry(k.lockDuration),
		redsync.SetTries(k.lockTries))

	return m, m.Lock()
}

func (k *keeper) decideCacheTTL(c Item) float64 {
	if c.GetTTLFloat64() > 0 {
		return c.GetTTLFloat64()
	}

	return k.defaultTTL.Seconds()
}

func (k *keeper) getCachedItem(key string, setStat bool) (value interface{}, err error) {
	client := k.connPool.Get()
	defer client.Close()

	if !setStat || !k.enableStat {
		return client.Do("GET", key)
	}

	client.Send("MULTI")
	client.Send("GET", key)
	client.Send("INCR", k.statPrefix+key)
	r, err := redigo.Values(client.Do("EXEC"))
	if err != nil {
		return
	}

	return r[0], nil
}

func (k *keeper) isLocked(key string) bool {
	client := k.lockConnPool.Get()
	defer client.Close()

	reply, err := client.Do("GET", "lock:"+key)
	if err != nil || reply == nil {
		return false
	}

	return true
}
