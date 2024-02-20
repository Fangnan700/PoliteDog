package pool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	startTime := time.Now().UnixMilli()

	var wg sync.WaitGroup
	wg.Add(5)

	p, _ := NewPool(5, 2)

	p.Submit(func() {
		fmt.Println("1111111111111111")
		time.Sleep(5 * time.Second)
		wg.Done()
	})

	p.Submit(func() {
		fmt.Println("2222222222222222")
		time.Sleep(5 * time.Second)
		wg.Done()
	})

	time.Sleep(5 * time.Second)

	p.Submit(func() {
		fmt.Println("3333333333333333")
		time.Sleep(5 * time.Second)
		wg.Done()
	})

	p.Submit(func() {
		fmt.Println("4444444444444444")
		time.Sleep(5 * time.Second)
		wg.Done()
	})

	time.Sleep(5 * time.Second)

	p.Submit(func() {
		fmt.Println("5555555555555555")
		time.Sleep(5 * time.Second)
		wg.Done()
	})

	wg.Wait()

	endTime := time.Now().UnixMilli()

	fmt.Println((endTime-startTime)/1000, " s")
}
