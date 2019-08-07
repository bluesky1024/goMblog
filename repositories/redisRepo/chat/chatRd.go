package chatRdRepo

import (
	"github.com/go-redis/redis"
)

type ChatRbRepository struct {
	redisPool *redis.ClusterClient
}

var logType = "chatRdRepo"

func NewChatRdRepo(redisPool *redis.ClusterClient) *ChatRbRepository {
	return &ChatRbRepository{
		redisPool: redisPool,
	}
}
