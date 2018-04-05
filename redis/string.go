package redis

import (
	redigo "github.com/garyburd/redigo/redis"
)

// Delete string value
func (r Redis) Delete(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("DEL", key))
}

// Get string value
func (r Redis) Get(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("GET", key))
}

// Set key and value
func (r Redis) Set(key, value string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("SET", key, value))
}

// HSet key and value
func (r Redis) HSet(key, key2 string, value string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("HSET", key, key2, value))
}

// HGet key
func (r Redis) HGet(key string, key2 string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("HGET", key, key2))

}

// SetWithNX with NX params
func (r Redis) SetWithNX(key, value string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("SET", key, value, "NX", "EX", expire))
}

// SetNX key and value
func (r Redis) SetNX(key, value string, expire int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("SETNX", key, value))
}

// SetEX key and value
func (r Redis) SetEX(key, value string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("SETEX", key, expire, value))
}
