package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "8.135.120.218"
	redisPort = "8010"
)

func init() {
	pool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 创建链接
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
		},
	}
}

// NewClient 创建Redis
func NewClient() redis.Conn {
	return pool.Get()
}
