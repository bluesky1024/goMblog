package redisSource

import (
	"github.com/bluesky1024/goMblog/config"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var feedMInstance *redis.Pool
var feedSInstance *redis.Pool
var feedLock *sync.Mutex = &sync.Mutex{}

func LoadFeedRdSour(master bool) (*redis.Pool, error) {
	var err error = nil
	if master {
		if feedMInstance == nil {
			feedLock.Lock()
			defer feedLock.Unlock()
			if feedMInstance == nil {
				redisConfig := conf.InitConfig("redisConfig.server")
				feedMInstance, err = LoadRedisSource(redisConfig["m_host"],
					redisConfig["m_port"])
				if err != nil {
					return nil, err
				}
			}
		}
		return feedMInstance, err
	} else {
		if feedSInstance == nil {
			feedLock.Lock()
			defer feedLock.Unlock()
			if feedMInstance == nil {
				redisConfig := conf.InitConfig("redisConfig.server")
				feedSInstance, err = LoadRedisSource(redisConfig["s_host"],
					redisConfig["s_port"])
				if err != nil {
					return nil, err
				}
			}
		}
		return feedSInstance, err
	}
}
