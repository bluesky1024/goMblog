package redisSource

import (
	"github.com/go-redis/redis"
	"sync"
)

var userInstance *redis.ClusterClient
var userLock *sync.Mutex = &sync.Mutex{}

func LoadUserRdSour() (*redis.ClusterClient, error) {
	var err error = nil
	if userInstance == nil {
		userLock.Lock()
		defer userLock.Unlock()
		if userInstance == nil {
			addrs := []string{
				"127.0.0.1:10011",
				"127.0.0.1:10012",
				"127.0.0.1:10013",
				"127.0.0.1:10014",
				"127.0.0.1:10015",
				"127.0.0.1:10016",
			}
			userInstance, err = LoadRedisSource(addrs)
			if err != nil {
				return nil, err
			}
		}
	}
	return userInstance, err
}
