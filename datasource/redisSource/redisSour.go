package redisSource

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func LoadRedisSource(host string, port string) (pool *redis.Pool, err error) {
	server := host + ":" + port
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", password); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			//if _, err := c.Do("SELECT", db); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			return c, nil
		},
	}
	return pool, err
}
