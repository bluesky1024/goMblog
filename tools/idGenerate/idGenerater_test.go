package idGenerate

import (
	"fmt"
	"sync"
	"testing"
)

func init() {
	err := InitUidPool(100)
	if err != nil {
		fmt.Println("init pool fail", err.Error())
		panic(err)
	}
}

func TestGenUidId(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_, err := GenUidId()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func TestGenUidIdWithMultiCoRoutine(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mid, err := GenUidId()
			fmt.Println(mid, err)
			if err != nil {
				fmt.Println(i, err.Error())
			}
		}(i)
	}
	wg.Wait()
}

func TestGenUidId_V2(t *testing.T) {

}
