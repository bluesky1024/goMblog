package userRdRepo

import (
	"github.com/gomodule/redigo/redis"
)

type UserRbRepository struct {
	RedisPoolM *redis.Pool
	RedisPoolS *redis.Pool
}

var logType = "userRdRepo"

func NewUserRdRepo(redisPoolM *redis.Pool, redisPoolS *redis.Pool) *UserRbRepository {
	return &UserRbRepository{
		RedisPoolM: redisPoolM,
		RedisPoolS: redisPoolS,
	}
}
