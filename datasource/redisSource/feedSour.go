package redisSource

import (
	"github.com/go-redis/redis"
	"sync"
)

var feedInstance *redis.ClusterClient
var feedLock *sync.Mutex = &sync.Mutex{}

func LoadFeedRdSour() (*redis.ClusterClient, error) {
	var err error = nil
	if feedInstance == nil {
		feedLock.Lock()
		defer feedLock.Unlock()
		if feedInstance == nil {
			addrs := []string{
				"127.0.0.1:10011",
				"127.0.0.1:10012",
				"127.0.0.1:10013",
				"127.0.0.1:10014",
				"127.0.0.1:10015",
				"127.0.0.1:10016",
			}
			feedInstance, err = LoadRedisClusterSource(addrs)
			if err != nil {
				return nil, err
			}
		}
	}
	return feedInstance, err
}
