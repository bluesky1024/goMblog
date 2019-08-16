package newIdGen

import (
	"fmt"
	"sync"
	"testing"
)

func init() {
	InitMidPool(100)
}

func TestGenMidId(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := GenMidId()
			//fmt.Println("res", mid, err)
			if err != nil {
				fmt.Println(i, err)
			}
		}(i)
	}
	wg.Wait()
	midPool.workerPool.Close()
}
