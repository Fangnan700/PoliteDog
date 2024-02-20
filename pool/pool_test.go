package pool

import (
	"math"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	KiB = 1024
	MiB = 1048567
)

const (
	Param    = 100
	PoolSize = 1000
	TestSize = 10000
	n        = 100000
)

const (
	RunTimes   = 1000000
	BenchParam = 10
)

var curMem uint64

func DemoFuc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			DemoFuc()
			wg.Done()
		}()
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)

	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("Memory usage: %d MB", curMem)
}

func TestWithPool(t *testing.T) {
	pool, _ := NewPool(math.MaxInt64, 10)
	defer pool.Release()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = pool.Submit(func() {
			DemoFuc()
			wg.Done()
		})
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)

	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("Memory usage: %d MB", curMem)
}
