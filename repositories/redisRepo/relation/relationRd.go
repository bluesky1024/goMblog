package relationRdRepo

import (
	"github.com/go-redis/redis"
)

type RelationRbRepository struct {
	redisPool *redis.ClusterClient
}

var logType = "relationRdRepo"

func NewRelationRdRepo(redisPool *redis.ClusterClient) *RelationRbRepository {
	return &RelationRbRepository{
		redisPool: redisPool,
	}
}
