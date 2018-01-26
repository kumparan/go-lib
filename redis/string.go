package redis

import (
	redigo "github.com/garyburd/redigo/redis"
)

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
