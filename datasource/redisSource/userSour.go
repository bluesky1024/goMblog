package redisSource

import (
	"github.com/bluesky1024/goMblog/config"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var userMInstance *redis.Pool
var userSInstance *redis.Pool
var userLock *sync.Mutex = &sync.Mutex{}

func LoadUserRdSour(master bool) (*redis.Pool, error) {
	var err error = nil
	if master {
		if userMInstance == nil {
			userLock.Lock()
			defer userLock.Unlock()
			if feedMInstance == nil {
				redisConfig := conf.InitConfig("redisConfig.server")
				userMInstance, err = LoadRedisSource(redisConfig["m_host"],
					redisConfig["m_port"])
				if err != nil {
					return nil, err
				}
			}
		}
		return userMInstance, err
	} else {
		if userSInstance == nil {
			userLock.Lock()
			defer userLock.Unlock()
			if userSInstance == nil {
				redisConfig := conf.InitConfig("redisConfig.server")
				userSInstance, err = LoadRedisSource(redisConfig["s_host"],
					redisConfig["s_port"])
				if err != nil {
					return nil, err
				}
			}
		}
		return userSInstance, err
	}
}
