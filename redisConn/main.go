package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gomodule/redigo/redis"
)

func RedisPool() (rPool redis.Pool){
	return redis.Pool{
		MaxIdle: 50,
		MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Fatal("ERROR: fail initializing the redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}
var redisPool = RedisPool()
var conn = redisPool.Get()

func ping() (error) {
    res, err := redis.String(conn.Do("PING"))
    if err != nil {
        return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	fmt.Println("PING response from Redis was %s", res)
	return err
}

func set(key string, value []byte) error {
	_, err := conn.Do("SET", key, value)
    if err != nil {
        return fmt.Errorf("error setting key %s to %s: %v", key, value, err)
	}
	return err
}

func get(key string) ([]byte, error) {
    var data []byte
    data, err := redis.Bytes(conn.Do("GET", key))
    if err != nil {
        return data, fmt.Errorf("error getting key %s: %v", key, err)
    }
    return data, err
}


func del(key string) error {
	_, err := conn.Do("DEL", key)
	if err != nil {
        return fmt.Errorf("error Deleting key %s: %v", key, err)
	}
	return err
}

func exists(key string) (bool, error) {
    ok, err := redis.Bool(conn.Do("EXISTS", key))
    if err != nil {
        return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
    }
    return ok, err
}

func main() {
	fmt.Println(ping())
}
