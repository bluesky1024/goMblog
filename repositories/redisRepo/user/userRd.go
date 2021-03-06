package userRdRepo

import (
	"github.com/go-redis/redis"
)

type UserRbRepository struct {
	redisPool *redis.ClusterClient
}

var logType = "userRdRepo"

func NewUserRdRepo(redisPool *redis.ClusterClient) *UserRbRepository {
	return &UserRbRepository{
		redisPool: redisPool,
	}
}
