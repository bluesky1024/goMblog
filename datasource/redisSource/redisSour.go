package redisSource

import (
	"github.com/go-redis/redis"
)

//redis集群
func LoadRedisClusterSource(addrs []string) (pool *redis.ClusterClient, err error) {
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	_, err = redisCluster.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisCluster, nil
}

//单点redis
func LoadRedisClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
