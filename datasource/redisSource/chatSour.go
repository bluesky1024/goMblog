package redisSource

import (
	"github.com/go-redis/redis"
	"sync"
)

var chatInstance *redis.ClusterClient
var chatLock = &sync.Mutex{}

func LoadChatRdSour() (*redis.ClusterClient, error) {
	var err error = nil
	if chatInstance == nil {
		chatLock.Lock()
		defer chatLock.Unlock()
		if chatInstance == nil {
			addrs := []string{
				"127.0.0.1:10011",
				"127.0.0.1:10012",
				"127.0.0.1:10013",
				"127.0.0.1:10014",
				"127.0.0.1:10015",
				"127.0.0.1:10016",
			}
			chatInstance, err = LoadRedisClusterSource(addrs)
			if err != nil {
				return nil, err
			}
		}
	}
	return chatInstance, err
}
