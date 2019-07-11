package relationDbRepo

import (
	"fmt"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"runtime"
	"sync"
	"testing"
)

var relationRepo *RelationDbRepository

func init() {
	relationSourceM, err := ds.LoadRelation(true)
	if err != nil {
		return
	}
	relationSourceS, err := ds.LoadRelation(false)
	if err != nil {
		return
	}
	relationRepo = NewRelationRepository(relationSourceM, relationSourceS)
}

func TestRelationDbRepository_UpdateFanCnt(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				res := relationRepo.UpdateFanCntByUid(123456, j)
				if !res {
					fmt.Println(res)
				}
			}
		}()
	}
	wg.Wait()
}

func TestRelationDbRepository_UpdateFollowCnt(t *testing.T) {
	res := relationRepo.UpdateFollowCntByUid(123, 10)
	fmt.Println(res)
}

func TestRelationDbRepository_GetFromInvalidTable(t *testing.T) {
	res := relationRepo.getRelationCntByUidsFromOneTable([]int64{1, 2, 3}, "123")
	fmt.Println(res)
}
