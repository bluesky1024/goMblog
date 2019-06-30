package redisSource

import (
	"github.com/go-redis/redis"
	"sync"
)

var relationInstance *redis.ClusterClient
var relationLock *sync.Mutex = &sync.Mutex{}

func LoadRelationRdSour() (*redis.ClusterClient, error) {
	var err error = nil
	if relationInstance == nil {
		userLock.Lock()
		defer userLock.Unlock()
		if relationInstance == nil {
			addrs := []string{
				"127.0.0.1:10011",
				"127.0.0.1:10012",
				"127.0.0.1:10013",
				"127.0.0.1:10014",
				"127.0.0.1:10015",
				"127.0.0.1:10016",
			}
			relationInstance, err = LoadRedisClusterSource(addrs)
			if err != nil {
				return nil, err
			}
		}
	}
	return relationInstance, err
}
