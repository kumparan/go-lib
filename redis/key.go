package redis

import (
	redigo "github.com/garyburd/redigo/redis"
)

// Expire set expiration time for a key
func (r Redis) Expire(key string, expiry int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("EXPIRE", key, expiry))
}

// Exists set expiration time for a key
func (r Redis) Exists(key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("EXISTS", key))
}

// Sort ordered the value lexicographically
func (r Redis) Sort(key string, alpha bool, asc bool) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	args := make([]interface{}, 3)
	args[0] = key
	if alpha {
		args[1] = "ALPHA"
	}
	if asc {
		args[2] = "ASC"
	} else {
		args[2] = "DESC"
	}

	return redigo.Strings(conn.Do("SORT", args...))
}
