package app

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

func openRedisPool() (*redis.Pool, error) {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				os.Getenv("REDIS_ADDR"),
				redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			// Ping the connection to ensure it's still alive
			_, err := c.Do("PING")
			return err
		},
		IdleTimeout: 240 * time.Second,
	}

	// Test Redis connection by getting and immediately releasing a connection
	conn := redisPool.Get()
	defer conn.Close() // Ensure the connection is returned to the pool
	_, err := conn.Do("PING")
	if err != nil {
		return redisPool, err
	}
	return redisPool, nil
}
