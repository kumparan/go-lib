package redis

import (
	redigo "github.com/garyburd/redigo/redis"
)

// RPush append values to the key
func (r Redis) RPush(key string, values []string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	args := make([]interface{}, len(values)+1)
	args[0] = key
	for i, value := range values {
		args[i+1] = value
	}
	return redigo.Int(conn.Do("RPUSH", args...))
}

// RPushEx append values to the key and set expiration for the key
func (r Redis) RPushEx(key string, values []string, expiry int) (int, error) {
	length, err := r.RPush(key, values)
	if err != nil {
		return length, err
	}

	return r.Expire(key, expiry)
}
