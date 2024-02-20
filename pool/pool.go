package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type signal struct{}

type Pool struct {
	// 最大运行worker数量
	cap int32

	// 正在运行worker数量
	running int32

	// 空闲worker的过期时间
	expire time.Duration

	// pool销毁信号
	release chan signal

	// 资源互斥锁
	lock sync.Mutex

	// 限制单次操作
	once sync.Once

	// 空闲worker队列
	workers []*worker
}

func NewPool(cap int32, expire int64) (*Pool, error) {
	if cap < 0 {
		return nil, errors.New("cap can not be less than 0")
	}
	if expire < 0 {
		return nil, errors.New("expire can not be less than 0")
	}

	pool := &Pool{}
	pool.cap = cap
	pool.expire = time.Duration(expire) * time.Second
	pool.release = make(chan signal, 1)

	go pool.expireWorker()
	return pool, nil
}

func (p *Pool) Submit(task func()) error {
	if len(p.release) > 0 {
		return errors.New("pool has been destroyed")
	}

	w := p.getWorker()
	w.task <- task
	w.pool.incRunning()

	return nil
}

func (p *Pool) Release() {
	p.once.Do(func() {
		p.lock.Lock()
		workers := p.workers
		for i, w := range workers {
			w.pool = nil
			w.task = nil
			workers[i] = nil
		}
		p.workers = nil
		p.lock.Unlock()

		p.release <- signal{}
	})
}

func (p *Pool) expireWorker() {
	ticker := time.NewTicker(p.expire)
	for range ticker.C {
		fmt.Printf("%v\n", p.workers)

		p.lock.Lock()
		var index int
		freeWorkers := p.workers
		if len(freeWorkers) > 0 {
			for i, w := range freeWorkers {
				if time.Now().Sub(w.last) <= p.expire {
					break
				} else {
					w.task <- nil
					index = i
					break
				}
			}

			freeWorkers = freeWorkers[index+1:]
			p.workers = freeWorkers
		}
		p.lock.Unlock()
	}
}

func (p *Pool) getWorker() *worker {
	var n int
	var freeWorkers []*worker

	freeWorkers = p.workers
	n = len(freeWorkers)

	if n > 0 {
		p.lock.Lock()
		w := freeWorkers[0]
		freeWorkers = freeWorkers[1:]
		p.workers = freeWorkers
		p.lock.Unlock()

		return w
	}

	if p.running < p.cap {
		w := &worker{
			pool: p,
			task: make(chan func(), 1),
		}
		w.run()

		return w
	} else {
		for {
			p.lock.Lock()
			freeWorkers = p.workers
			n = len(freeWorkers)

			if n <= 0 {
				p.lock.Unlock()
				continue
			}

			w := freeWorkers[0]
			freeWorkers = freeWorkers[1:]
			p.workers = freeWorkers
			p.lock.Unlock()

			return w
		}
	}
}

func (p *Pool) putWorker(w *worker) {
	w.last = time.Now()
	p.lock.Lock()
	p.workers = append(p.workers, w)
	p.lock.Unlock()
}

func (p *Pool) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

func (p *Pool) decRunning() {
	atomic.AddInt32(&p.running, -1)
}
