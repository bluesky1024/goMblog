package redisLock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSingleRedisLock(t *testing.T) {
	fmt.Println("test redis lock")
	InitRedisLockSrv("127.0.0.1:6379")
	lockSrv := LoadSrv()

	tempLock, err := lockSrv.TryGetLock("testLock", 10*time.Second, 10*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("do something")

	fmt.Println("lock param:", tempLock)
	//tempLock.UnLock()
}

func TestMultiRedisLock(t *testing.T) {
	InitRedisLockSrv("127.0.0.1:6379")
	lockSrv := LoadSrv()

	lockName := "multiLock"

	wg := sync.WaitGroup{}
	for i := 0; i < 260; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tempLock, err := lockSrv.TryGetLock(lockName, 10*time.Second, 100*time.Millisecond)
			if err != nil {
				fmt.Println(i, err.Error())
				return
			}
			//fmt.Println(i, "get lock success", "do something ...")
			tempLock.UnLock()
			//fmt.Println(i, "end")
		}(i)
	}
	wg.Wait()
	fmt.Println("finish test")
}
