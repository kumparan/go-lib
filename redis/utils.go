package redis

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// Ping nodoc
func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redigo.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

// Pong nodoc
func Pong() {
	client := Pool.Get()
	defer client.Close()
	n, err := client.Do("INFO")
	if err != nil {
		log.Println("ERROR:" + err.Error())
	}
	log.Printf("info=%s", n)
}

// Execute nodoc
func Execute(command string, key string, value string) interface{} {

	var content interface{}
	var err error

	client := Pool.Get()
	defer client.Close()

	if value != "" {
		content, err = redigo.Values(client.Do(command, key, value))
	} else {
		content, err = redigo.Values(client.Do(command, key))
	}

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR:" + err.Error())
	}

	return content
}

// Expire nodoc
func Expire(key string, redisTTL int64) {
	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("EXPIRE", key, redisTTL)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR EXPIRE: " + key + "; " + err.Error())
	}
}

// Set nodoc
func Set(key string, value string, redisTTL int64) {
	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("SET", key, value)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR SET: " + key + "; " + err.Error())
	}

	Expire(key, redisTTL)
}

// Get nodoc
func Get(key string) string {
	client := Pool.Get()
	defer client.Close()

	isExists, err := redigo.Int(client.Do("EXISTS", key))
	if err != nil && err != redis.ErrNil {
		log.Println("ERROR EXISTS: " + key + "; " + err.Error())
	}

	if isExists == 0 {
		// data not exists/expired from redis
		return ""
	}

	content, err := redigo.String(client.Do("GET", key))

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR GET: " + key + "; " + err.Error())
	}

	return content
}

// DeletePrefix nodoc
func DeletePrefix(key string) {
	client := Pool.Get()
	defer client.Close()

	keys, err := redis.Strings(client.Do("KEYS", key))

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR KEYS: " + key + "; " + err.Error())
	}

	for _, v := range keys {
		Delete(v)
	}

}

// Delete nodoc
func Delete(key string) {

	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("DEL", key)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR DEL: " + key + "; " + err.Error())
	}

}

// HSet nodoc
func HSet(key string, key2 string, value string, redisTTL int64) {
	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("HSET",
		key, key2, value)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR HSET: " + key + "; " + key2 + "; " + err.Error())
	}

	Expire(key, redisTTL)
}

// HGet nodoc
func HGet(key string, key2 string) string {
	client := Pool.Get()
	defer client.Close()

	n, err := client.Do("HGET", key, key2)
	if err != nil && err != redis.ErrNil {
		log.Println("ERROR HGET1: " + key + "; " + key2 + "; " + err.Error())
	}

	if n != nil {
		content, err := redigo.String(client.Do("HGET", key, key2))

		if err != nil && err != redis.ErrNil {
			log.Println("ERROR HGET2: " + key + "; " + key2 + "; " + err.Error())
		}

		return content
	} else {
		return "nil"
	}
}

// SAdd nodoc
func SAdd(key string, value string, redisTTL int64) {
	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("SADD", key, value)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR SADD: " + key + "; " + err.Error())
	}

	if redisTTL > 0 {
		Expire(key, redisTTL)
	}
}

// SMembers nodoc
func SMembers(key string) interface{} {
	client := Pool.Get()
	defer client.Close()

	content, err := redigo.Values(client.Do("SMEMBERS", key))

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR SMEMBERS: " + key + "; " + err.Error())
	}

	return content
}

// SISMember nodoc
func SISMember(key string, value string) int {
	client := Pool.Get()
	defer client.Close()

	content, err := redigo.Int(client.Do("SISMEMBER", key, value))

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR SISMEMBER: " + key + "; " + err.Error())
	}

	return content
}

// SRem nodoc
func SRem(key string, value string) {
	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("SREM", key, value)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR SREM: " + key + "; " + err.Error())
	}
}

// RPUSH :nodoc:
func RPUSH(key string, value string) {
	client := Pool.Get()
	defer client.Close()

	_, err := client.Do("RPUSH", key, value)

	if err != nil && err != redis.ErrNil {
		log.Println("ERROR RPUSH: " + key + "; " + err.Error())
	}
}
