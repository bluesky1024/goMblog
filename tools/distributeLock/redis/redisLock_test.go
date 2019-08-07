package redisLock

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	once sync.Once
)

func testSetup() {
	addrs := []string{
		"127.0.0.1:10011",
		"127.0.0.1:10012",
		"127.0.0.1:10013",
		"127.0.0.1:10014",
		"127.0.0.1:10015",
		"127.0.0.1:10016",
	}
	InitRedisLockSrv(addrs)
}

func TestGetRoutineId(t *testing.T) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	fmt.Println(id)
}

func TestSingleRedisLock(t *testing.T) {
	fmt.Println("test redis lock")
	once.Do(testSetup)
	lockSrv := LoadSrv()

	tempLock, err := lockSrv.TryGetLock("multiLock", 1000*time.Second, 1000*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("do something")

	fmt.Println("lock param:", tempLock)

	err = tempLock.UnLock()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("do something")

	fmt.Println("lock param:", tempLock)
	//tempLock.UnLock()
}

func TestMultiRedisLock(t *testing.T) {
	once.Do(testSetup)
	lockSrv := LoadSrv()

	lockName := "multiLock"

	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tempLock, err := lockSrv.TryGetLock(lockName, 10*time.Second, 100*time.Millisecond)
			if err != nil {
				fmt.Println(i, err)
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

func TestLockWithExtend(t *testing.T) {
	once.Do(testSetup)
	fmt.Println("test redis lock")
	lockSrv := LoadSrv()

	tempLock, err := lockSrv.TryGetLock("testLock", 1*time.Second, 2*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("do something")

	ctx := context.Background()
	ctxCancel, cancel := context.WithCancel(ctx)
	tempLock.ExtendLock(ctxCancel)

	time.Sleep(10 * time.Second)

	cancel()
	tempLock.UnLock()
}

func TestGoCoroutineA(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < 10000; i++ {
		temp := 0
		for j := 0; j < 1000; j++ {
			temp = temp * (j + i)
		}
	}
	wg.Done()
	wg.Wait()
}

func TestGoCoroutineB(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			temp := 0
			for j := 0; j < 10000; j++ {
				temp = temp * (j + i)
			}
		}(i)
	}
	wg.Wait()
}
