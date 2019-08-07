package feedRdRepo

import (
	"fmt"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"testing"
)

func TestFeedRbRepository_GetFeeds(t *testing.T) {
	redisCluster, err := redisSource.LoadFeedRdSour()
	if err != nil {
		panic(err)
	}
	repo := NewFeedRdRepo(redisCluster)
	mids, err := repo.GetFeeds(1, 0, 0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(mids)
}

func TestFeedRbRepository_GetFirstMid(t *testing.T) {
	redisCluster, err := redisSource.LoadFeedRdSour()
	if err != nil {
		panic(err)
	}
	repo := NewFeedRdRepo(redisCluster)
	mid, err := repo.GetFirstMid(1, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(mid)
}

func TestFeedRbRepository_GetLastMid(t *testing.T) {
	redisCluster, err := redisSource.LoadFeedRdSour()
	if err != nil {
		panic(err)
	}
	repo := NewFeedRdRepo(redisCluster)
	mid, err := repo.GetLastMid(1, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(mid)
}

func TestFeedRbRepository_RemoveMids(t *testing.T) {
	redisCluster, err := redisSource.LoadFeedRdSour()
	if err != nil {
		panic(err)
	}
	repo := NewFeedRdRepo(redisCluster)
	mids := []int64{111, 222}
	err = repo.RemoveMids(1, 0, mids)
	if err != nil {
		panic(err)
	}
}
