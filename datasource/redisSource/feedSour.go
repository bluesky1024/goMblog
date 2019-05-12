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
				"192.168.0.60:10011",
				"192.168.0.60:10012",
				"192.168.0.60:10013",
				"192.168.0.60:10014",
				"192.168.0.60:10015",
				"192.168.0.60:10016",
			}
			feedInstance, err = LoadRedisSource(addrs)
			if err != nil {
				return nil, err
			}
		}
	}
	return feedInstance, err
}
