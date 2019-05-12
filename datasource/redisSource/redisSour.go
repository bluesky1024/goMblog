package redisSource

import (
	"github.com/go-redis/redis"
)

func LoadRedisSource(addrs []string) (pool *redis.ClusterClient, err error) {
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	pong, err := redisCluster.Ping().Result()
	if err != nil || pong != "PONG" {
		return nil, err
	}
	return redisCluster, nil
}
